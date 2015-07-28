package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"

	"time"
)

const (
	ProjectsTable = "project"
	CommentsTable = "comment"
)

func QueryProjects() orm.QuerySeter {
	return orm.NewOrm().QueryTable(ProjectsTable)
}

func QueryComments() orm.QuerySeter {
	return orm.NewOrm().QueryTable(CommentsTable)
}

func CheckProjectExists(name string) bool {
	return QueryProjects().Filter("Name", name).Exist()
}

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
			q = q.SetCond(orm.NewCondition().And("Status__Name__in", value))
		case "technos":
			q = q.SetCond(orm.NewCondition().And("Technos__Skill__Name__in", value))
		case "themes":
			q = q.SetCond(orm.NewCondition().And("Themes__Theme__Name__in", value))
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
		if _, err = o.LoadRelated(project, "Technos"); err != nil {
			return
		}
		if _, err = o.LoadRelated(project, "Themes"); err != nil {
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
	if err := QueryProjects().SetCond(orm.NewCondition().Or("Id", nameOrId).Or("Name", nameOrId)).RelatedSel().One(project); err != nil {
		return nil, err
	}
	return loadProjectInfo(project)
}

func GetProjectByID(id int) (*models.Project, error) {
	project := &models.Project{}
	if err := QueryProjects().Filter("Id", id).RelatedSel().One(project); err != nil {
		return nil, err
	}
	return loadProjectInfo(project)
}

func AddAndGetProject(project *models.Project) (*models.Project, error) {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, err
	}
	status, err := GetProjectStatusByName(project.StatusName)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	project.Status = status
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
	if err := setProjectRelation(project, historyItem); err != nil {
		o.Rollback()
		return nil, err
	}
	project.Id = int(id)
	o.Commit()
	return GetProjectByID(int(id))
}

func AddCommentToProject(comment *models.Comment, project *models.Project) error {
	o := orm.NewOrm()
	id, err := o.Insert(comment)
	if err != nil {
		return err
	}
	comment.Id = int(id)
	_, err = o.QueryM2M(project, "Comments").Add(comment)
	return err
}

func EditComment(comment *models.Comment) error {
	_, err := QueryComments().Filter("Id", comment.Id).Update(orm.Params{
		"updated": time.Now(),
		"message": comment.Message,
	})
	return err
}

func GetCommentByID(id int) (*models.Comment, error) {
	comment := &models.Comment{}
	err := QueryComments().Filter("Id", id).RelatedSel().One(comment)
	return comment, err
}

func loadProjectInfo(project *models.Project) (*models.Project, error) {
	o := orm.NewOrm()
	if _, err := o.LoadRelated(project, "Members"); err != nil {
		return nil, err
	}
	if err := loadEveryInfoOfUsers(project.Members); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "History"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "Themes"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "Technos"); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(project, "Comments"); err != nil {
		return nil, err
	}
	for _, c := range project.Comments {
		if _, err := o.LoadRelated(c, "Author"); err != nil {
			return nil, err
		}
	}
	return project, nil
}

func EditAndGetProject(project *models.Project) (*models.Project, error) {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, err
	}
	status, err := GetProjectStatusByName(project.StatusName)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	project.Status = status
	now := time.Now()
	params := orm.Params{
		"name":                project.Name,
		"shortDescription":    project.ShortDescription,
		"completeDescription": project.CompleteDescription,
		"repository":          project.Repository,
		"website":             project.Website,
		"updated":             now,
		"status_id":           status.Id,
	}
	if project.Manager != nil {
		params["manager_id"] = project.Manager.Id
	}
	if _, err := o.QueryTable(ProjectsTable).Filter("Id", project.Id).Update(params); err != nil {
		o.Rollback()
		return nil, err
	}
	historyItem, err := AddAndGetHistoryEvent(models.HistoryItemTypeCreated, project)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	if err := clearProjectRelation(project); err != nil {
		o.Rollback()
		return nil, err
	}
	if err := setProjectRelation(project, historyItem); err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return GetProjectByID(project.Id)
}

func setProjectRelation(project *models.Project, historyItem *models.HistoryItem) error {
	o := orm.NewOrm()
	if len(project.Members) != 0 {
		if _, err := o.QueryM2M(project, "Members").Add(project.Members); err != nil {
			return err
		}
	}
	if len(project.Themes) != 0 {
		if _, err := o.QueryM2M(project, "Themes").Add(project.Themes); err != nil {
			return err
		}
	}
	if len(project.Technos) != 0 {
		if _, err := o.QueryM2M(project, "Technos").Add(project.Technos); err != nil {
			return err
		}
	}
	if _, err := o.QueryM2M(project, "History").Add(historyItem); err != nil {
		return err
	}
	return nil
}

func clearProjectRelation(project *models.Project) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	o := orm.NewOrm()
	o.QueryM2M(project, "Members").Clear()
	o.QueryM2M(project, "Themes").Clear()
	o.QueryM2M(project, "Technos").Clear()
	o.QueryM2M(project, "History").Clear()
	return err
}
