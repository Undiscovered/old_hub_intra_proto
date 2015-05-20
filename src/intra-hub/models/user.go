package models

import (
	"github.com/astaxie/beego/orm"
    "github.com/astaxie/beego/validation"
)

func init() {
	orm.RegisterModel(&User{})
}

type User struct {
	Id        int
	Login     string `orm:"unique" form:"login"`
	FirstName string
	LastName  string
	Email     string
	Picture   string
	Password  string `form:"password"`
	Promotion string
	City      string
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Login"},
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
