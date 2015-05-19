package main

import (
	_ "intra-hub/routers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
    _ "intra-hub/models"
)

func main() {
    beego.SetStaticPath("/css", "static/css")
    beego.SetStaticPath("/js", "static/js")
    beego.SetStaticPath("/img", "static/img")

    orm.RegisterDriver("mysql", orm.DR_MySQL)
    orm.RegisterDataBase("default", "mysql", "root@/orm_test?charset=utf8")
    orm.RunCommand()
	beego.Run()

}

