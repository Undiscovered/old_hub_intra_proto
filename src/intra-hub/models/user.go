package models

import (
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
    "github.com/beego/i18n"
    "encoding/json"
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
}

type User struct {
	Id          int        `json:"id"`
	Login       string     `json:"login" orm:"unique;size(128)" form:"login"`
	FirstName   string     `json:"firstName" orm:"size(128)"`
	LastName    string     `json:"lastName" orm:"size(128)"`
	Email       string     `json:"email" orm:"size(128)" form:"email"`
	Picture     string     `json:"picture" orm:"size(128)"`
	Password    string     `json:"password" orm:"size(128)" form:"password"`
	PhoneNumber string     `json:"phoneNumber" orm:"size(16)" form:"phoneNumber"`
	Token       string     `json:"token,omitempty" orm:"size(128)"`
	Tech4Derog  bool       `json:"tech4Derog" form:"tech4Derog"`
	Promotion   *Promotion `json:"promotion" orm:"null;rel(fk)"`
	City        *City      `json:"city" orm:"null;rel(fk)"`
	Group       *Group     `json:"group" orm:"null;rel(fk)"`
	Projects    []*Project `json:"projects" orm:"rel(m2m)"`
	Skills      []*Skill   `json:"skills" orm:"rel(m2m);rel_through(intra-hub/models.UserSkill)"`
}

func (u *User) Clean() *User {
    u.Password = ""
    return u
}

func (u *User) Name() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Login", "FirstName", "LastName"},
	}
}

func (u *User) Values() []string {
	return []string{u.Login, u.FirstName, u.LastName, u.Email, u.Picture,
		u.Password, strconv.FormatInt(int64(u.Group.Id), 10),
		strconv.FormatInt(int64(u.Promotion.Id), 10), strconv.FormatInt(int64(u.City.Id), 10)}
}

func (u *User) Valid(v *validation.Validation) {
	if len(u.Login) > 15 {
		v.SetError("Login", "invalid login")
	}
	if len(u.Password) == 0 {
		v.SetError("Password", "empty password")
	}
}

func (u *User) ToJSON(locale string) (string, error) {
    u.City.Name = i18n.Tr(locale, u.City.Name)
    u.Group.Name = i18n.Tr(locale, u.Group.Name)
    js, err := json.Marshal(u)
    if err != nil {
        return "", err
    }
    return string(js), nil
}

func GetUserFields() string {
	return "login, first_name, last_name, email, picture, password, group_id, promotion_id, city_id"
}
