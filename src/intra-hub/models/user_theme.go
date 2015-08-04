package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&UserThemes{})
}

type UserThemes struct {
	Id    int
	User  *User  `orm:"rel(fk)"`
	Theme *Theme `orm:"rel(fk)"`
	Level int
}

func (u *UserThemes) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "User"},
	}
}
