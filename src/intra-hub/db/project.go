package db

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
)

const (
	ProjectsTable = "project"
)

func GetProjectsPaginated(page, limit int, queryFilter map[string]interface{}) (itemPaginated *models.ItemPaginated, err error) {
	projects := make([]*models.Project, 0)
	o := orm.NewOrm()
	q := o.QueryTable(ProjectsTable)
	page -= 1
	for key, value := range queryFilter {
		switch value.(type) {
		case []string:
			if value.([]string)[0] == "" {
				continue
			}
		case string:
			if value.(string) == "" {
				continue
			}
		}
		beego.Warn(key, value)
		switch key {
		case "promotions":
			q = q.SetCond(orm.NewCondition().And("Members__User__Promotion__Name__in", value))
		case "cities":
			q = q.SetCond(orm.NewCondition().And("Members__User__City__Name__in", value))
		case "student":
			q = q.SetCond(orm.NewCondition().And("Members__User__Login__icontains", value).
				Or("Members__User__FirstName__icontains", value).Or("Members__User__LastName__icontains", value))
		case "managers":
			q = q.SetCond(orm.NewCondition().And("Manager__Login__in", value))
		case "status":
			q = q.SetCond(orm.NewCondition().And("Status__in", value))
		case "name":
			q = q.SetCond(orm.NewCondition().And("Name__icontains", value))
		}
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
    if len(project.Themes) != 0 {
        if _, err := o.QueryM2M(project, "Themes").Add(project.Themes); err != nil {
            o.Rollback()

        }
    }
	if _, err := o.QueryM2M(project, "History").Add(historyItem); err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return GetProjectByID(int(id))
}
