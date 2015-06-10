package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
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

func (c *ProjectController) IntroView() {
	c.TplNames = "project/intro.html"
}

func (c *ProjectController) ListView() {
	c.TplNames = "project/list.html"
	handleError := func(err error) {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
	}
	queryFilter := make(map[string]interface{})
	queryFilter["promotions"] = strings.Split(c.GetString("promotions", ""), ",")
	queryFilter["cities"] = strings.Split(c.GetString("cities", ""), ",")
	queryFilter["managers"] = strings.Split(c.GetString("managers", ""), ",")
	queryFilter["status"] = strings.Split(c.GetString("status", ""), ",")
	queryFilter["student"] = c.GetString("student", "")
	queryFilter["name"] = c.GetString("name", "")
	page, err := c.GetInt("page")
	if err != nil {
		handleError(err)
		return
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		handleError(err)
		return
	}
	if page <= 0 {
		c.Redirect(fmt.Sprintf("/projects?page=1&limit=%d", limit), 301)
		return
	}
	if limit == 0 {
		limit = 25
	}
	paginatedItems, err := db.GetProjectsPaginated(page, limit, queryFilter)
	if err != nil {
		handleError(err)
		return
	}
	paginatedItems.SetPagesToShow()
	promotions, err := db.GetEveryPromotion()
	if err != nil {
		handleError(err)
		return
	}
	cities, err := db.GetEveryCities()
	if err != nil {
		handleError(err)
		return
	}
	managers, err := db.GetManagers()
	if err != nil {
		handleError(err)
		return
	}
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
	c.Data["Cities"] = cities
	c.Data["Promotions"] = promotions
	c.Data["Limit"] = limit
	c.Data["PaginatedItems"] = paginatedItems
	c.Data["HasNextPage"] = paginatedItems.CurrentPage+1 <= paginatedItems.TotalPageCount
	c.Data["HasPreviousPage"] = paginatedItems.CurrentPage != 1
	c.Data["ShowGoToFirst"] = paginatedItems.PagesToShow[0] != 1
	c.Data["ShowGoToLast"] = paginatedItems.PagesToShow[len(paginatedItems.PagesToShow)-1] != paginatedItems.TotalPageCount
}

func (c *ProjectController) SingleView() {
	c.TplNames = "project/single.html"
	project, err := db.GetProjectByIDOrName(c.GetString(":nameOrId"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/projects/list?page=1&limit=15", 301)
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
	themes, err := db.GetEveryThemes()
	if err != nil {
		beego.Error(err)
		c.flash.Data["error"] = err.Error()
		return
	}
	c.Data["Themes"] = themes
	c.Data["Status"] = models.EveryProjectStatus
	c.Data["Managers"] = managers
}

func (c *ProjectController) Add() {
	c.TplNames = "project/add.html"
	project := &models.Project{}
	if err := c.ParseForm(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	valid := validation.Validation{}
	if b, err := valid.Valid(project); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	} else if !b {
		beego.Error(valid.Errors[0])
		c.SetErrorAndRedirect(fmt.Errorf(valid.Errors[0].String()))
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
	beego.Warn(project.Themes, project.ThemesID)
	project, err := db.AddAndGetProject(project)
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/projects/"+strconv.FormatInt(int64(project.Id), 10), 301)
}
