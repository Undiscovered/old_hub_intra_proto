package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"intra-hub/db"
	"intra-hub/models"
	"strconv"
)

type ProjectController struct {
	BaseController
}

func (c *ProjectController) NestedPrepare() {
//	if !c.isLogged {
//		beego.Warn("Not Logged")
//		c.Redirect("/", 301)
//		return
//	}
}

func (c *ProjectController) ListView() {
	offset, err := c.GetInt("offset")
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	if limit == 0 {
		limit = 25
	}
	projects, err := db.GetProjectsPaginated(offset, limit)
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Projects"] = projects
	c.TplNames = "project/list.html"
}

func (c *ProjectController) SingleView() {
	c.TplNames = "project/single.html"
	id, err := strconv.Atoi(c.Ctx.Input.Params[":id"])
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	project, err := db.GetProjectByID(id)
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Project"] = project
}

func (c *ProjectController) AddView() {
	c.TplNames = "project/add.html"
	managers, err := db.GetManagers()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
}

func (c *ProjectController) Add() {
	c.TplNames = "project/add.html"
	project := &models.Project{}
	c.ParseForm(project)
	beego.Warning(project)
	valid := validation.Validation{}
	if b, err := valid.Valid(project); err != nil {
		beego.Error(err)
        c.SetErrorAndRedirect("/projects/add", err)
		return
	} else if !b {
		beego.Error(valid.Errors[0])
		c.flash.Data["error"] = valid.Errors[0].String()
		c.flash.Store(&c.Controller)
		c.Redirect("/projects/add", 303)
		return
	}
	if project.ManagerLogin != "--" {
		manager, err := db.GetUserByLogin(project.ManagerLogin)
		if err != nil {
			beego.Error(err)
			c.flash.Data["error"] = err.Error()
			c.flash.Store(&c.Controller)
			c.Redirect("/projects/add", 303)
		}
		project.Manager = manager
	}
	project, err := db.AddProject(project)
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		c.flash.Store(&c.Controller)
		c.Redirect("/projects/add", 303)
		return
	}
	c.Redirect("/projects/"+strconv.FormatInt(int64(project.Id), 10), 301)
}
