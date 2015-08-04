package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&Theme{})
}

type Theme struct {
	Id    int    `json:"id"`
	Name  string `json:"name" orm:"unique;size(128)"`
	Level int    `json:"level" orm:"-"`
}
