package db
import (
    "intra-hub/models"
    "github.com/astaxie/beego/orm"
)

const (
    ThemeTable = "theme"
)

func QueryThemes() orm.QuerySeter {
    return orm.NewOrm().QueryTable(ThemeTable)
}

func GetAllThemes() (themes []*models.Theme, err error) {
    themes = make([]*models.Theme, 0)
    _, err = QueryThemes().All(&themes)
    return
}

func AddTheme(theme *models.Theme) error {
    _, err := orm.NewOrm().Insert(theme)
    return err
}