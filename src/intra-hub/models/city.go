package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&City{})
}

type City struct {
	Id   int `json:"id"`
	Name string `json:"name" orm:"unique;size(128)"`
    LocalizedName string `json:"localizedName" orm:"-"`
}
