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

func GetGroupByName(name string) (*models.Group, error) {
	group := &models.Group{}
	return group, QueryGroup().Filter("Name", name).One(group)
}

func GetEveryGroups() ([]*models.Group, error) {
	groups := make([]*models.Group, 0)
	_, err := QueryGroup().All(&groups)
	return groups, err
}
