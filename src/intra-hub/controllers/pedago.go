package controllers

import (
	"github.com/astaxie/beego"
	"github.com/saschpe/tribool"
	"intra-hub/db"
)

type PedagoController struct {
	BaseController
}

func (c *PedagoController) NestedPrepare() {
	c.RequirePedago()
}

func (c *PedagoController) ValidateProjectView() {
	c.TplNames = "pedago/project-validation.html"
	validationStr := c.GetString(":validation")
	var validation tribool.Tribool
	switch validationStr {
	case "indeterminate":
		validation = tribool.Indeterminate
	case "refused":
		validation = tribool.False
	case "validated":
		validation = tribool.True
	}
	userProjects, err := db.GetEveryUserProjectsByValidation(validation)
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
		return
	}
	c.Data["UserProjects"] = userProjects
}

func (c *PedagoController) ValidateProject() {
	c.EnableRender = false
	defer c.Redirect("/pedago/validation/indeterminate", 301)
	userID, err := c.GetInt(":userId")
	if err != nil {
		c.flash.Data["error"] = err.Error()
		return
	}
	projectID, err := c.GetInt(":projectId")
	if err != nil {
		c.flash.Data["error"] = err.Error()
		return
	}
	validation, err := c.GetInt(":validation")
	if err != nil {
		c.flash.Data["error"] = err.Error()
		return
	}
	if err := db.ValidatePedagogicallyUser(userID, projectID, tribool.Tribool(validation)); err != nil {
		c.flash.Data["error"] = err.Error()
		return
	}
}
