package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&City{})
}

type City struct {
	Id   int
	Name string `orm:"unique;size(128)"`
}
