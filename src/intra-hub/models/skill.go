package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&Skill{})
}

type Skill struct {
	Id    int    `json:"id"`
	Name  string `json:"name" orm:"unique;size(128)"`
	Level int    `json:"level"`
}
