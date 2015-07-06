package controllers

import (
	"github.com/astaxie/beego"
	"intra-hub/db"
)

type AdminController struct {
	BaseController
}

func (c *AdminController) NestedPrepare() {
	c.RequireManager()
}

func (c *AdminController) Get() {
	c.TplNames = "admin/layout.html"
	themes, err := db.GetEveryThemes()
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	skills, err := db.GetEverySkills()
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Data["Themes"] = themes
	c.Data["Skills"] = skills
}
