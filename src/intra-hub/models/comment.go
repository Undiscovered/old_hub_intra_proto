package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(&Comment{})
}

type Comment struct {
	Id      int       `form:"id"`
	Message string    `json:"message" form:"message" orm:"null;type(text)"`
	Author  *User     `orm:"rel(fk)"`
	Created time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated time.Time `json:"updated" orm:"auto_now;type(datetime)"`

	// Non persistent fields

	ProjectID int `form:"projectId" orm:"-"`
}
