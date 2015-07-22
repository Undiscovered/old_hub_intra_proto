package controllers

import (
	"encoding/json"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type SkillController struct {
	BaseController
}

func (c *SkillController) Post() {
	if !c.user.IsManager() {
		jsonErr := simplejson.New()
		jsonErr.Set("error", "forbidden")
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	skill := &models.Skill{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), skill); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	skill, err := db.AddAndGetSkill(skill)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	c.Data["json"] = skill
	c.ServeJson()
}

func (c *SkillController) Delete() {
	if !c.user.IsManager() {
		jsonErr := simplejson.New()
		jsonErr.Set("error", "forbidden")
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	skillID, err := c.GetInt(":id", -1)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	if err := db.DeleteSkillByID(skillID); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	c.Data["json"] = jsonOK
	c.ServeJson()
}
