package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(&Group{})
}

type Group struct {
	Id   int
	Name string `orm:"unique;size(128)"`
}
