package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	GroupTable = "group"
)

func QueryGroup() orm.QuerySeter {
	return orm.NewOrm().QueryTable(GroupTable)
}

func GetGroupByNames(name string) (*models.Group, error) {
	group := &models.Group{}
	return group, QueryGroup().Filter("Name__exact", name).One(group)
}
