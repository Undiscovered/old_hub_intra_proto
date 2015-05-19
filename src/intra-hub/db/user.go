package db
import "github.com/astaxie/beego/orm"

const (
    UserTable = "user"
)

func QueryUser() orm.QuerySeter {
    return orm.NewOrm().QueryTable(UserTable)
}