package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	HistoryItemTypeCreated = "CREATED"
	HistoryItemTypeUpdated = "UPDATED"
)

func init() {
	orm.RegisterModel(&HistoryItem{})
}

type HistoryItem struct {
	Id      int
	Type    string
	Date    time.Time `orm:"auto_now_add;type(datetime)"`
	Project *Project  `orm:"rel(fk)"`
}
