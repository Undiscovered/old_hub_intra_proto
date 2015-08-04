package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&UserSkills{})
}

type UserSkills struct {
	Id    int
	User  *User  `orm:"rel(fk)"`
	Skill *Skill `orm:"rel(fk)"`
	Level int
}

func (u *UserSkills) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "User"},
	}
}
