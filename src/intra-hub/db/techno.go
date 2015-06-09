package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	TechnoTable = "techno"
)

func QueryTechnos() orm.QuerySeter {
	return orm.NewOrm().QueryTable(TechnoTable)
}

func GetAllTechnos() (technos []*models.Techno, err error) {
	technos = make([]*models.Techno, 0)
	_, err = QueryTechnos().All(&technos)
	return
}

func AddAndGetTechno(techno *models.Techno) (*models.Techno, error) {
	id, err := orm.NewOrm().Insert(techno)
	if err != nil {
		return nil, err
	}
	techno.Id = int(id)
	return techno, nil
}

func DeleteTechnoByID(id int) error {
	_, err := orm.NewOrm().Delete(&models.Techno{Id: id})
	return err
}
