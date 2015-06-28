package models

import (
    "time"
    "github.com/astaxie/beego/orm"
)

func init() {
    orm.RegisterModel(&Comment{})
}

type Comment struct {
	Id      int
	Message string    `json:"message"`
	Created time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated time.Time `json:"updated" orm:"auto_now;type(datetime)"`
}
