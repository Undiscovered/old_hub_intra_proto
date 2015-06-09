package controllers

import (
	"encoding/json"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type TechnoController struct {
	BaseController
}

func (c *TechnoController) Post() {
	techno := &models.Techno{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), techno); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	techno, err := db.AddAndGetTechno(techno)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	c.Data["json"] = techno
	c.ServeJson()
}

func (c *TechnoController) Delete() {
	technoID, err := c.GetInt(":id", -1)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	if err := db.DeleteTechnoByID(technoID); err != nil {
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
