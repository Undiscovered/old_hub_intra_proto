package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(&UserSkill{})
}

type UserSkill struct {
	Id    int
	User  *User  `orm:"rel(fk)"`
	Skill *Skill `orm:"rel(fk)"`
	Level int
}