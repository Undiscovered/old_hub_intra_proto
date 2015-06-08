package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
)

const (
	UserGroupStudent = "STUDENT"
	UserGroupManager = "MANAGER"
	UserGroupAdmin   = "ADMIN"
	UserGroupPedago  = "PEDAGO"
)

var (
	EveryUserGroups = []string{UserGroupStudent, UserGroupPedago, UserGroupManager, UserGroupAdmin}
)

func init() {
	orm.RegisterModel(&User{})
	orm.RegisterModel(&City{})
	orm.RegisterModel(&Promotion{})
}

type City struct {
	Id   int
	Name string `orm:"unique;size(128)"`
}

type Promotion struct {
	Id   int
	Name string `orm:"unique;size(128)"`
}

type User struct {
	Id        int        `json:"id"`
	Login     string     `json:"login" orm:"unique;size(128)" form:"login"`
	FirstName string     `json:"firstName" orm:"size(128)"`
	LastName  string     `json:"lastName" orm:"size(128)"`
	Email     string     `json:"email" orm:"size(128)"`
	Picture   string     `json:"picture" orm:"size(128)"`
	Password  string     `json:"password" orm:"size(128)" form:"password"`
	Promotion *Promotion `json:"promotion" orm:"null;rel(fk)"`
	City      *City      `json:"city" orm:"null;rel(fk)"`
	Groups    []*Group   `json:"groups" orm:"rel(m2m)"`
	Projects  []*Project `json:"projects" orm:"rel(m2m)"`
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Login", "FirstName", "LastName"},
	}
}

func (u *User) Values() []string {
	return []string{u.Login, u.FirstName, u.LastName, u.Email, u.Picture,
		u.Password, strconv.FormatInt(int64(u.Promotion.Id), 10), strconv.FormatInt(int64(u.City.Id), 10)}
}

func (u *User) Valid(v *validation.Validation) {
	if len(u.Login) > 15 {
		v.SetError("Login", "invalid login")
	}
	if len(u.Password) == 0 {
		v.SetError("Password", "empty password")
	}
}

func GetUserFields() string {
	return "login, first_name, last_name, email, picture, password, promotion_id, city_id"
}
