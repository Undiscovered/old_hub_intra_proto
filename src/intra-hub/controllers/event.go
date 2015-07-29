package controllers

import (
	_ "strconv"
	_ "strings"

	"intra-hub/db"
	"intra-hub/models"

	_ "github.com/astaxie/beego"
	_ "github.com/astaxie/beego/validation"
)

type EventController struct {
	BaseController
}

func (c *EventController) CalendarAdd() {
	c.EnableRender = false
	c.RequireAdmin()

	cal := &models.Calendar{}
	if err := c.ParseForm(cal); err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	if _, err := db.AddAndGetCalendar(cal); err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	c.Redirect("/admin", 301)
}

func (c *EventController) CalendarAdminView() {
	c.RequireAdmin()
	calendarId, err := c.GetInt(":id", -1)
	if err != nil || calendarId != -1 {
		c.TplNames = "calendar/admin-one.html"
		calendar, err := db.GetCalendarByID(calendarId)
		if err == nil {
			c.Data["Calendar"] = calendar
		}
	}
}

func (c *EventController) CalendarsView() {
	calendarId, err := c.GetInt(":id", -1)
	if err != nil || calendarId != -1 {
		c.TplNames = "calendar/one.html"
		calendar, err := db.GetCalendarByID(calendarId)
		if c.user == nil || calendar.Public == false && !c.user.IsAdmin() {
			c.Redirect("/home", 301)
			return
		}
		if err == nil {
			c.Data["Calendar"] = calendar
		}
	} else {
		c.TplNames = "calendar/index.html"
		calendars, err := db.GetEveryCalendars()
		if err == nil {
			c.Data["Calendars"] = calendars
		}
	}
}
