package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&Techno{})
}

type Techno struct {
	Id   int    `json:"id"`
	Name string `json:"name" orm:"unique;size(128)"`
}
