package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(&User{})
}

type User struct {
	Id        int
	Login     string `orm:"unique"`
	FirstName string
	LastName  string
	Email     string
	Picture   string
	Password  string
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

func GetUserFields() string {
    return "login, first_name, last_name, email, picture, password, promotion, city"
}