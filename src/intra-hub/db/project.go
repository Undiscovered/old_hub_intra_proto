package db
import (
    "intra-hub/models"
    "github.com/astaxie/beego/orm"
)

const (
    ProjectsTable = "project"
)

func QueryProjects() orm.QuerySeter {
    return orm.NewOrm().QueryTable(ProjectsTable)
}

func GetProjectsPaginated(offset, limit int) (projects []*models.Project, err error) {
    projects = make([]*models.Project, 0)
    if _, err = QueryProjects().Offset(offset).Limit(limit).All(&projects); err != nil {
        return
    }
    return
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
    if _, err := o.LoadRelated(project, "History"); err != nil {
        return nil, err
    }
    return project, nil
}

func AddProject(project *models.Project) (*models.Project, error) {
    o := orm.NewOrm()
    if err := o.Begin(); err != nil {
        return nil, err
    }
    id, err :=  o.Insert(project)
    if err != nil {
        o.Rollback()
        return nil, err
    }
    historyItem, err := addHistoryEvent(models.HistoryItemTypeCreated, project)
    if err != nil {
        o.Rollback()
        return nil, err
    }
    project.Id = int(id)
    if _, err := o.QueryM2M(project, "Members").Add(project.Members); err != nil {
        o.Rollback()
        return nil, err
    }
    if _, err := o.QueryM2M(project, "History").Add(historyItem); err != nil {
        o.Rollback()
        return nil, err
    }
    o.Commit()
    return GetProjectByID(int(id))
}