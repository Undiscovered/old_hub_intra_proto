package main

import (
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    _ "intra-hub/routers"
    _ "github.com/go-sql-driver/mysql"
    _ "intra-hub/models"
    _ "intra-hub/tasks"
    "time"
    "os"
)

const (
    driverSQL = "mysql"
    aliasDbName = "default"
    databaseName = "intra_hub"
    username = "root"
    password = ""
    maxIdleConns = 30
    maxOpenConns = 30
    optionsDatabaseConnections = "?charset=utf8"
)

func main() {
    // Set session on
    beego.SessionOn = true

    beego.EnableAdmin = true

    // Serve static files
    beego.SetStaticPath("/css", "static/css")
    beego.SetStaticPath("/js", "static/js")
    beego.SetStaticPath("/img", "static/img")

    // Set logger
    os.Create("logs/test.log")
    beego.SetLogger("file", `{"filename":"logs/test.log"}`)
    beego.SetLogFuncCall(true)
    beego.SetLevel(beego.LevelInformational)

    // Set the ORM parameters
    orm.RegisterDriver(driverSQL, orm.DR_MySQL)
    orm.RegisterDataBase(aliasDbName, driverSQL, username + password + "@/" + databaseName + optionsDatabaseConnections)
    orm.SetMaxIdleConns(aliasDbName, maxIdleConns)
    orm.SetMaxOpenConns(aliasDbName, maxOpenConns)
    orm.DefaultTimeLoc = time.UTC
    orm.RunCommand()

    // Run the app
	beego.Run()

}

