package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/saschpe/tribool"
)

func init() {
	orm.RegisterModel(&UserProjects{})
}

type UserProjects struct {
	Id                     int
	User                   *User           `orm:"rel(fk)"`
	Project                *Project        `orm:"rel(fk)"`
	PedagogicallyValidated tribool.Tribool `orm:"default(1)"`
}

func (u *UserProjects) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Project"},
	}
}
