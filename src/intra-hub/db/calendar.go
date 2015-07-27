package db

import (
	"github.com/astaxie/beego/orm"
	"intra-hub/models"
)

const (
	CalendarCollection = "calendar"
)

func QueryCalendar() orm.QuerySeter {
	return orm.NewOrm().QueryTable(CalendarCollection)
}

func AddCalendar(calendar *models.Calendar) error {
	_, err := orm.NewOrm().Insert(calendar)
	return err
}

func GetEveryCalendars() ([]*models.Calendar, error) {
	calendars := make([]*models.Calendar, 0)
	_, err := QueryCalendar().All(&calendars)
	return calendars, err
}

func DeleteCalendar(id int) error {
	_, err := QueryCalendar().Filter("Id", id).Delete()
	return err
}
