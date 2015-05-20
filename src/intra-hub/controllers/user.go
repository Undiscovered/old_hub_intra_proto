package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"intra-hub/db"
	"intra-hub/models"
)

type UserController struct {
	BaseController
}

// TODO: API Mode
func (c *UserController) Login() {
	user := &models.User{}

	if c.apiMode {

	} else {
        flash := beego.ReadFromRequest(&c.Controller)
        c.TplNames = "login.html"
		c.ParseForm(user)
		beego.Warning(user)
		valid := validation.Validation{}
		if b, err := valid.Valid(user); err != nil {
			beego.Error(err)
			flash.Data["error"] = err.Error()
			return
		} else if !b {
            beego.Error(valid.Errors)
            flash.Data["error"] = valid.Errors[0].String()
            return
        }
		user, err := db.CheckUserCredentials(user)
		if err != nil {
			beego.Error(err)
			flash.Data["error"] = err.Error()
			return
		}
		c.Data["User"] = user
        c.TplNames = "index.html"
	}
}

func (c *UserController) LoginView() {
	c.TplNames = "login.html"
}
