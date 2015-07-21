package main

import (
	"strings"
	"time"

	"intra-hub/confperso"
	"intra-hub/db"
	_ "intra-hub/models"
	_ "intra-hub/routers"
	_ "intra-hub/tasks"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	"github.com/beego/i18n"
	"github.com/eknkc/dateformat"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strconv"
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

	dateFormat = "dddd DD MMMM YYYY HH:mm:ss"
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

	// Set session on
	beego.SessionOn = true
	beego.SessionProvider = driverSQL
	beego.SessionSavePath = username + ":" + password + "@/" + databaseName + optionsDatabaseConnections

	beego.EnableAdmin = true

	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"

	// Serve static files
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/img", "static/img")

	// Set the ORM parameters
	orm.RegisterDriver(driverSQL, orm.DR_MySQL)
	orm.RegisterDataBase(aliasDbName, driverSQL, username+":"+password+"@/"+databaseName+optionsDatabaseConnections)
	orm.SetMaxIdleConns(aliasDbName, maxIdleConns)
	orm.SetMaxOpenConns(aliasDbName, maxOpenConns)
	orm.DefaultTimeLoc = time.UTC
	orm.RunCommand()
	db.PopulateDatabase()

	// Add Default templating functions
	beego.AddFuncMap("i18n", i18n.Tr)
	incr := func(arg int) string {
		return strconv.FormatInt(int64(arg+1), 10)
	}
	decr := func(arg int) string {
		return strconv.FormatInt(int64(arg-1), 10)
	}
	randomizeLabel := func() string {
		labels := []string{"success", "warning", "danger", "info", "primary", "default"}
		return labels[rand.Intn(len(labels))]
	}
	toJSON := func(val interface{}) string {
		js, _ := json.Marshal(val)
		return string(js)
	}
	datefr := func(val time.Time) string {
		return dateformat.FormatLocale(val, dateFormat, dateformat.French)
	}
	beego.AddFuncMap("incr", incr)
	beego.AddFuncMap("decr", decr)
	beego.AddFuncMap("randLabel", randomizeLabel)
	beego.AddFuncMap("toJSON", toJSON)
	beego.AddFuncMap("datefr", datefr)

}

func main() {
	// Run the app
	beego.Run()
}
