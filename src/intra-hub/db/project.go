package db
import (
    "intra-hub/models"
    "github.com/astaxie/beego/orm"
    "github.com/astaxie/beego"
)

const (
    ProjectsTable = "project"
)

func QueryProjects() orm.QuerySeter {
    return orm.NewOrm().QueryTable(ProjectsTable)
}

func GetProjectsPaginated(page, limit int) (itemPaginated *models.ItemPaginated, err error) {
    projects := make([]*models.Project, 0)
    o := orm.NewOrm()
    q := o.QueryTable(ProjectsTable)
    page -= 1
    if _, err = q.Offset(page * limit).Limit(limit).RelatedSel().All(&projects); err != nil {
        return
    }
    for _, project := range projects {
        if _, err := o.LoadRelated(project, "Members"); err != nil {
            return nil, err
        }
    }
    count, err := q.Count()
    if err != nil {
        return
    }
    itemPaginated = &models.ItemPaginated{
        Items: projects,
        ItemCount: len(projects),
        TotalItemCount: int(count),
        CurrentPage: page + 1,
        TotalPageCount: int(count) / limit + 1,
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