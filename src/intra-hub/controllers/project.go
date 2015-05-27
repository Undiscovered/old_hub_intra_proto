package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"intra-hub/db"
	"intra-hub/models"
	"strconv"
    "fmt"
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
    c.TplNames = "project/list.html"
	page, err := c.GetInt("page")
	if err != nil {
		beego.Error(err)
        c.SetErrorAndRedirect(err)
		return
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		beego.Error(err)
        c.SetErrorAndRedirect(err)
		return
	}
    if page <= 0 {
        c.Redirect(fmt.Sprintf("/projects?page=1&limit=%d", limit), 301)
        return
    }
	if limit == 0 {
		limit = 25
	}
    beego.Warn(beego.BeeTemplates)
	paginatedItems, err := db.GetProjectsPaginated(page, limit)
	if err != nil {
		beego.Error(err)
        c.SetErrorAndRedirect(err)
		return
	}
    c.Data["Limit"] = limit
	c.Data["PaginatedItems"] = paginatedItems
    c.Data["HasNextPage"] = paginatedItems.CurrentPage != paginatedItems.TotalPageCount
    c.Data["HasPreviousPage"] = paginatedItems.CurrentPage != 1
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
		c.SetErrorAndRedirect(err)
		return
	} else if !b {
		beego.Error(valid.Errors[0])
        c.SetErrorAndRedirect(err)
		return
	}
	if project.ManagerLogin != "--" {
		manager, err := db.GetUserByLogin(project.ManagerLogin)
		if err != nil {
			beego.Error(err)
            c.SetErrorAndRedirect(err)
            return
		}
		project.Manager = manager
	}
	project, err := db.AddProject(project)
	if err != nil {
		beego.Error(err)
        c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/projects/"+strconv.FormatInt(int64(project.Id), 10), 301)
}
