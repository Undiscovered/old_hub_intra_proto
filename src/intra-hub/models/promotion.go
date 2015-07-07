package models

import "github.com/astaxie/beego/orm"

const (
	ExternalPromotion = "external"
)

func init() {
	orm.RegisterModel(&Promotion{})
}

type Promotion struct {
	Id   int    `json:"id"`
	Name string `json:"name" orm:"unique;size(128)"`
}
