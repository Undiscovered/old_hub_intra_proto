package models
import "github.com/astaxie/beego/orm"

func init() {
    orm.RegisterModel(&Techno{})
}

type Techno struct {
    Id int
	Name string `orm:"unique;size(128)"`
}
