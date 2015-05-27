package db

import (
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
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
