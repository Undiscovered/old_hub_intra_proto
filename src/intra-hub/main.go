package main

import (
	"os"
	"strings"
	"time"

	"intra-hub/confperso"
	"intra-hub/db"
	_ "intra-hub/models"
	_ "intra-hub/routers"
	_ "intra-hub/tasks"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	_ "github.com/go-sql-driver/mysql"
)

const (
	driverSQL                  = "mysql"
	aliasDbName                = confperso.AliasDbName
	databaseName               = confperso.DatabaseName
	username                   = confperso.Username
	password                   = confperso.Password
	maxIdleConns               = 150
	maxOpenConns               = 150
	optionsDatabaseConnections = "?charset=utf8"
)

// langType represents a language type.
type langType struct {
	Lang, Name string
}

func init() {
	// This is just used to get every languages files (see app.conf files)
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	langTypes := make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	// Then we load every language files with i18n.SetMessage
	for _, lang := range langs {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}

	orm.Debug = true

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

	// Set the ORM parameters
	orm.RegisterDriver(driverSQL, orm.DR_MySQL)
	orm.RegisterDataBase(aliasDbName, driverSQL, username+":"+password+"@/"+databaseName+optionsDatabaseConnections)
	orm.SetMaxIdleConns(aliasDbName, maxIdleConns)
	orm.SetMaxOpenConns(aliasDbName, maxOpenConns)
	orm.DefaultTimeLoc = time.UTC
	orm.RunCommand()
	db.PopulateDatabase()
}

func main() {
	// Run the app
	beego.Run()
}
