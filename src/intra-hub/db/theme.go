package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
	"intra-hub/services/cache"
)

const (
	ThemeTable = "theme"
)

func QueryThemes() orm.QuerySeter {
	return orm.NewOrm().QueryTable(ThemeTable)
}

func GetEveryThemes() (themes []*models.Theme, err error) {
	themes = make([]*models.Theme, 0)
	_, err = QueryThemes().OrderBy("Name").All(&themes)
	go cache.SetThemes(themes)
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

func EditAndGetTheme(theme *models.Theme) (*models.Theme, error) {
	_, err := orm.NewOrm().Update(theme)
	return theme, err
}

func DeleteThemeByID(id int) error {
	_, err := orm.NewOrm().Delete(&models.Theme{Id: id})
	return err
}
