package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&Promotion{})
}

type Promotion struct {
	Id   int
	Name string `orm:"unique;size(128)"`
}
