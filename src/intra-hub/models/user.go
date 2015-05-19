package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(&User{})
}

type User struct {
	Id    int
	Login string
}
