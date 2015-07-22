package controllers

import (
	_"strconv"
	_"strings"
	
	"intra-hub/db"
	"intra-hub/models"

	_"github.com/astaxie/beego"
	_"github.com/astaxie/beego/validation"
)

type CalendarController struct {
	BaseController
}

func (c *CalendarController) CalendarAdd() {
	c.EnableRender = false
	c.RequireAdmin()

	cal := &models.Calendar{}
	err := c.ParseForm(cal);
	_, err = db.AddAndGetCalendar(cal)
	if err != nil {
		c.SetErrorAndRedirect(err)
	} else {
		c.Redirect("/admin", 301)
	}

}

func (c *CalendarController) CalendarAdminView() {
	c.RequireAdmin()
	calendarId, err := c.GetInt(":id", -1)
	if err != nil || calendarId != -1 {
		c.TplNames = "calendar/admin-single.html"
		calendar, err := db.GetCalendarByID(calendarId)
		if err == nil {
			c.Data["Calendar"] = calendar
		}
	}
}

func (c *CalendarController) CalendarsView() {
	calendarId, err := c.GetInt(":id", -1)
	if err != nil || calendarId != -1 {
		c.TplNames = "calendar/single.html"
		calendar, err := db.GetCalendarByID(calendarId)
		if c.user == nil || calendar.Public == false && !c.user.IsAdmin() {
			c.Redirect("/home", 301)
			return;
		}
		if err == nil {
			c.Data["Calendar"] = calendar
		}
	} else {
		c.TplNames = "calendar/list.html"
		calendars, err := db.GetEveryCalendars()
		if err == nil {
			c.Data["Calendars"] = calendars
		}
	}
}
