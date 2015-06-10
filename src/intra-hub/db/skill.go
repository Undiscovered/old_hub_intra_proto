package db

import (
	"intra-hub/models"

	"github.com/astaxie/beego/orm"
)

const (
	SkillTable = "skill"
)

func QuerySkills() orm.QuerySeter {
	return orm.NewOrm().QueryTable(SkillTable)
}

func GetEverySkills() (skills []*models.Skill, err error) {
	skills = make([]*models.Skill, 0)
	_, err = QuerySkills().All(&skills)
	return
}

func AddAndGetSkill(skill *models.Skill) (*models.Skill, error) {
	id, err := orm.NewOrm().Insert(skill)
	if err != nil {
		return nil, err
	}
	skill.Id = int(id)
	return skill, nil
}

func DeleteSkillByID(id int) error {
	_, err := orm.NewOrm().Delete(&models.Skill{Id: id})
	return err
}
