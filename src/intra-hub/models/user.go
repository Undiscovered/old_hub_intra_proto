package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

const (
	UserGroupStudent = "STUDENT"
	UserGroupManager = "MANAGER"
	UserGroupAdmin   = "ADMIN"
)

var (
	EveryUserGroups = []string{UserGroupStudent, UserGroupManager, UserGroupAdmin}
)

func init() {
	orm.RegisterModel(&User{})
}

type User struct {
	Id        int        `json:"id"`
	Login     string     `json:"login" orm:"unique" form:"login"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Picture   string     `json:"picture"`
	Password  string     `json:"password" form:"password"`
	Promotion string     `json:"promotion"`
	City      string     `json:"city"`
	Groups    []*Group   `json:"groups" orm:"rel(m2m)"`
	Projects  []*Project `json:"projects" orm:"rel(m2m)"`
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Login", "FirstName", "LastName"},
	}
}

func (u *User) Values() []string {
	return []string{u.Login, u.FirstName, u.LastName, u.Email, u.Picture, u.Password, u.Promotion, u.City}
}

func (u *User) Valid(v *validation.Validation) {
	if len(u.Login) > 8 {
		v.SetError("Login", "invalid login")
	}
	if len(u.Password) == 0 {
		v.SetError("Password", "empty password")
	}
}

func GetUserFields() string {
	return "login, first_name, last_name, email, picture, password, promotion, city"
}
