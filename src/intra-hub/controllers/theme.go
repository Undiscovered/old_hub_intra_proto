package controllers

import (
	"encoding/json"

	"intra-hub/db"
	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type ThemeController struct {
	BaseController
}

func (c *ThemeController) Post() {
	theme := &models.Theme{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), theme); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	theme, err := db.AddAndGetTheme(theme)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	c.Data["json"] = theme
	c.ServeJson()
}

func (c *ThemeController) Delete() {
	themeID, err := c.GetInt(":id", -1)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	if err := db.DeleteThemeByID(themeID); err != nil {
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
