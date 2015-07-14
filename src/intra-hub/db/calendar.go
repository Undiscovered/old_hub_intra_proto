package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	CalendarTable = "calendar"
)

func QueryCalendars() orm.QuerySeter {
	return orm.NewOrm().QueryTable(CalendarTable)
}

func LoadRelated(i interface{}, field string) {
	o := orm.NewOrm()
	o.LoadRelated(i, field)
}

func GetCalendarByID(id int) (*models.Calendar, error) {
	var calendar models.Calendar
	err := QueryCalendars().Filter("id", id).One(&calendar)
	if err == nil {
		LoadRelated(&calendar, "Events")
	}
	return &calendar, err
}

func GetEveryCalendars() (calendars []*models.Calendar, err error) {
	calendars = make([]*models.Calendar, 0)
	_, err = QueryCalendars().All(&calendars)

	for i := 0; i < len(calendars); i++ {
		LoadRelated(calendars[i], "Events")
	}
	
	return calendars, err
}

func AddAndGetCalendar(calendar *models.Calendar) (*models.Calendar, error) {
	id, err := orm.NewOrm().Insert(calendar)
	if err != nil {
		return nil, err
	}
	calendar.Id = int(id)
	return calendar, nil
}

func DeleteCalendarByID(id int) error {
	_, err := orm.NewOrm().Delete(&models.Calendar{Id: id})
	return err
}
