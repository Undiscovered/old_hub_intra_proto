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

func AddAndGetTheme(theme *models.Theme) (*models.Theme, error) {
	id, err := orm.NewOrm().Insert(theme)
	if err != nil {
		return nil, err
	}
	theme.Id = int(id)
	return theme, nil
}

func DeleteThemeByID(id int) error {
	_, err := orm.NewOrm().Delete(&models.Theme{Id: id})
	return err
}
