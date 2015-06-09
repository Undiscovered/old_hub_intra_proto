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

func GetGroupsByNames(names ...string) (groups []*models.Group, err error) {
	_, err = QueryGroup().Filter("Name__in", names).All(&groups)
	return
}
