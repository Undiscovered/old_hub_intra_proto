package db

import (
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
)

// Populate the database with default values.
func PopulateDatabase() {
	// Add Groups
	{
		i, err := QueryGroup().PrepareInsert()
		if err != nil {
			panic(err)
		}
		for _, groupName := range models.EveryUserGroups {
			group := &models.Group{
				Name: groupName,
			}
			i.Insert(group)
		}
		if err := i.Close(); err != nil {
			panic(err)
		}
	}
	// Add Status
	{
		i, err := QueryProjectStatus().PrepareInsert()
		if err != nil {
			panic(err)
		}
		for _, statusName := range models.EveryProjectStatus {
			status := &models.ProjectStatus{
				Name: statusName,
			}
			i.Insert(status)
		}
		if err := i.Close(); err != nil {
			panic(err)
		}
	}
	// Add External Promotion
	{
		orm.NewOrm().Insert(&models.Promotion{Name: models.ExternalPromotion})
	}
}
