package db

import (
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
)

func addHistoryEvent(eventType string, project *models.Project) (item *models.HistoryItem, err error) {
	o := orm.NewOrm()
	item = &models.HistoryItem{
		Project: project,
		Type:    eventType,
	}
	id, err := o.Insert(item)
    item.Id = int(id)
	return
}