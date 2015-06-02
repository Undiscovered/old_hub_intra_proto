package db
import (
    "intra-hub/models"
    "github.com/astaxie/beego/orm"
)

const (
    CityTable = "city"
)

func QueryCities() orm.QuerySeter {
    return orm.NewOrm().QueryTable(CityTable)
}

func GetEveryCities() (cities []*models.City, err error) {
    _, err = QueryCities().All(&cities)
    return
}