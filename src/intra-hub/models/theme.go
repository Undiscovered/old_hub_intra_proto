package models
import "github.com/astaxie/beego/orm"

func init() {
    orm.RegisterModel(&Theme{})
}

type Theme struct {
    Id int
	Name string `orm:"unique;size(128)"`
}
