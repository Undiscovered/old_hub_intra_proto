package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	ProjectStatusTable = "project_status"
)

func QueryProjectStatus() orm.QuerySeter {
	return orm.NewOrm().QueryTable(ProjectStatusTable)
}

func GetProjectStatusByNames(names ...string) (groups []*models.ProjectStatus, err error) {
	_, err = QueryProjectStatus().Filter("Name__in", names).All(&groups)
	return
}

func GetProjectStatusByName(name string) (*models.ProjectStatus, error) {
	status := &models.ProjectStatus{}
	return status, QueryProjectStatus().Filter("Name__exact", name).One(status)
}
