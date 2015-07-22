package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
	"intra-hub/services/cache"
)

const (
	CityTable = "city"
)

func QueryCities() orm.QuerySeter {
	return orm.NewOrm().QueryTable(CityTable)
}

func GetEveryCities() (cities []*models.City, err error) {
	_, err = QueryCities().All(&cities)
	cache.SetCities(cities)
	return
}
