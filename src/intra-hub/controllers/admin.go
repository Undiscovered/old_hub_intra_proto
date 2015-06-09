package controllers

import (
	"github.com/astaxie/beego"
	"intra-hub/db"
)

type AdminController struct {
	BaseController
}

func (c *AdminController) Get() {
	c.TplNames = "admin/layout.html"
	themes, err := db.GetAllThemes()
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	skills, err := db.GetAllSkills()
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	technos, err := db.GetAllTechnos()
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Data["Themes"] = themes
	c.Data["Skills"] = skills
	c.Data["Technos"] = technos
}
