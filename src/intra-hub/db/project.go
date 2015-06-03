package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
	"log"
)

const (
	ProjectsTable = "project"
)

func C(page, limit int, promotions, cities []string, name string) (itemPaginated *models.ItemPaginated, err error) {
	projects := make([]*models.Project, 0)
	raw := "SELECT DISTINCT T0.`id`, T0.`name`, T0.`short_description`, T0.`status`, T0.`manager_id`, T0.`created`, T0.`updated`, " +
		"T1.`id`, T1.`login`, T1.`first_name`, T1.`last_name`, T1.`email`, T1.`picture`, T1.`password`, T1.`promotion_id`, T1.`city_id`," +
		"T2.`id`, T2.`name`," +
		"T3.`id`, T3.`name`" +
		"FROM `project` T0 " +
		"LEFT OUTER JOIN `user` T1 ON T1.`id` = T0.`manager_id` " +
		"LEFT OUTER JOIN `promotion` T2 ON T2.`id` = T1.`promotion_id` " +
		"LEFT OUTER JOIN `city` T3 ON T3.`id` = T1.`city_id` " +
		"INNER JOIN `user_projects` T4 ON T4.`project_id` = T0.`id` " +
		"INNER JOIN `user` T5 ON T5.`id` = T4.`user_id` " +
		"LEFT OUTER JOIN `city` T6 ON T6.`id` = T5.`city_id` "
	values := make([]interface{}, 0)
	if promotions[0] != "" {
		raw += "WHERE T2.`name` IN (?) "
		values = append(values, promotions)
	}
	if cities[0] != "" {
		raw += "WHERE T3.`name` IN (?) "
		values = append(values, cities)
	}
	if name != "" {
		raw += "WHERE T0.`name` COLLATE UTF8_GENERAL_CI LIKE (%?%) "
		values = append(values, name)
	}
	page -= 1
	raw += fmt.Sprintf("LIMIT %d OFFSET %d", limit, page*limit)
	beego.Warn(values...)
	var count int64
	if len(values) > 0 {
		if count, err = orm.NewOrm().Raw(raw, values...).QueryRows(&projects); err != nil {
			return
		}
	} else {
		if count, err = orm.NewOrm().Raw(raw).QueryRows(&projects); err != nil {
			return
		}
	}
	beego.Warn(count)
	itemPaginated = &models.ItemPaginated{
		Items:          projects,
		ItemCount:      len(projects),
		TotalItemCount: int(count),
		CurrentPage:    page + 1,
		TotalPageCount: int(count)/limit + 1,
	}
	log.Printf("%#v\n", projects[0])
	return
}

func GetProjectsPaginated(page, limit int, promotions, cities, managers, status []string, name string) (itemPaginated *models.ItemPaginated, err error) {
	//    return C(page, limit, promotions, cities, name)
	projects := make([]*models.Project, 0)
	o := orm.NewOrm()
	q := o.QueryTable(ProjectsTable)
	page -= 1
	if promotions[0] != "" {
		q = q.Filter("Members__User__Promotion__Name__in", promotions)
	}
	if cities[0] != "" {
		q = q.Filter("Members__User__City__Name__in", cities)
	}
	if managers[0] != "" {
		q = q.Filter("Manager__Login__in", managers)
	}
	if status[0] != "" {
		q = q.Filter("Status__in", status)
	}
	if name != "" {
		q = q.Filter("Name__icontains", name)
	}
	if _, err = q.RelatedSel().All(&projects); err != nil {
		return
	}
	m := make(map[int]bool)
	// Clear duplicates
	for i := len(projects) - 1; i >= 0; i-- {
		project := projects[i]
		if m[project.Id] {
			copy(projects[i:], projects[i+1:])
			projects[len(projects)-1] = nil
			projects = projects[:len(projects)-1]
		} else {
			m[project.Id] = true
		}
	}
	count := len(projects)
	tmp := page*limit + limit
	if tmp > len(projects) {
		tmp = len(projects)
	}
	projects = projects[page*limit : tmp]
	beego.Warn(projects, len(projects), cap(projects))
	for _, project := range projects {
		if _, err = o.LoadRelated(project, "Members"); err != nil {
			return
		}
		if err = loadEveryInfoOfUsers(project.Members); err != nil {
			return
		}
	}
	itemPaginated = &models.ItemPaginated{
		Items:          projects,
		ItemCount:      len(projects),
		TotalItemCount: count,
		CurrentPage:    page + 1,
		TotalPageCount: count/limit + 1,
	}
	return
}

func GetProjectByIDOrName(nameOrId string) (*models.Project, error) {
	project := &models.Project{}
	o := orm.NewOrm()
	q := o.QueryTable(ProjectsTable)
	beego.Warn(nameOrId)
	if err := q.SetCond(orm.NewCondition().Or("Id", nameOrId).Or("Name", nameOrId)).RelatedSel().One(project); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "Members"); err != nil {
		return nil, err
	}
	if err := loadEveryInfoOfUsers(project.Members); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "History"); err != nil {
		return nil, err
	}
	return project, nil
}

func GetProjectByID(id int) (*models.Project, error) {
	project := &models.Project{}
	o := orm.NewOrm()
	q := o.QueryTable(ProjectsTable)
	if err := q.Filter("Id", id).RelatedSel().One(project); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "Members"); err != nil {
		return nil, err
	}
	if err := loadEveryInfoOfUsers(project.Members); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "History"); err != nil {
		return nil, err
	}
	return project, nil
}

func AddAndGetProject(project *models.Project) (*models.Project, error) {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, err
	}
	id, err := o.Insert(project)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	historyItem, err := AddAndGetHistoryEvent(models.HistoryItemTypeCreated, project)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	project.Id = int(id)
	if len(project.Members) != 0 {
		if _, err := o.QueryM2M(project, "Members").Add(project.Members); err != nil {
			o.Rollback()
			return nil, err
		}
	}
	if _, err := o.QueryM2M(project, "History").Add(historyItem); err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return GetProjectByID(int(id))
}
