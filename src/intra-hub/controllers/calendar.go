package controllers

import (
	"github.com/astaxie/beego"
	"intra-hub/db"
	"intra-hub/models"
)

type CalendarController struct {
	BaseController
}

func (c *CalendarController) Add() {
	c.RequireManager()
	c.EnableRender = false
	calendar := &models.Calendar{}
	if err := c.ParseForm(calendar); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	beego.Warn(c.Input())
	beego.Warn(calendar)
	if err := db.AddCalendar(calendar); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/admin#calendars", 301)
}

func (c *CalendarController) Delete() {
	c.RequireManager()
	c.EnableRender = false
	id, _ := c.GetInt(":id", 0)
	if err := db.DeleteCalendar(id); err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/admin#calendars", 302)
}
