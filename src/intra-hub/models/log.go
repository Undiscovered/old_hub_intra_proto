package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(&Log{})
}

type Log struct {
	Id       int
	Action   string
	Table    string
	TargetID int
	User     *User     `orm:"rel(fk)"`
	Date     time.Time `orm:"auto_now_add;type(datetime)"`
}
