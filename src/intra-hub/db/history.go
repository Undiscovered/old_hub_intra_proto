package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

func AddAndGetHistoryEvent(eventType string, project *models.Project) (item *models.HistoryItem, err error) {
	o := orm.NewOrm()
	item = &models.HistoryItem{
		Project: project,
		Type:    eventType,
	}
	id, err := o.Insert(item)
	item.Id = int(id)
	return
}
