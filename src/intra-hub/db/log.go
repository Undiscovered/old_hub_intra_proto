package db

import (
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
)

func AddLog(user *models.User, action, table string, targetID int) error {
	l := &models.Log{
		Action:   action,
		Table:    table,
		TargetID: targetID,
		User:     user,
	}
	_, err := orm.NewOrm().Insert(l)
	return err
}
