package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(&Calendar{})
	orm.RegisterModel(&CalendarEvent{})
}

type Calendar struct {
	Id     int              `json:"id"`
	Name   string           `json:"name" orm:"size(128)" form:"name"`
	Public bool             `json:"public" form:"public"`
	Events []*CalendarEvent `json:"events" orm:"reverse(many)"`
}

type CalendarEvent struct {
	Id          int       `json:"id"`
	Title       string    `json:"title" orm:"size(128)" form:"title"`
	Description string    `json:"description" orm:"null;type(text)" form:"description"`
	PictureFile string    `json:"-" orm:"size(128)"`
	Location    string    `json:"location" orm:"size(128)" form:"location"`
	Begin       time.Time `json:"begin" form:"begin"`
	End         time.Time `json:"end" form:"end"`
	Calendar    *Calendar `json:"calendar" orm:"rel(fk)"`

	Picture string `json:"picture" orm:"-"`
}
