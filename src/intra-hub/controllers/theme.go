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
	if c.user == nil || !c.user.IsManager() {
		jsonErr := simplejson.New()
		jsonErr.Set("error", "forbidden")
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(403)
		c.ServeJson()
		return
	}
	theme := &models.Theme{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), theme); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	theme, err := db.AddAndGetTheme(theme)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	c.Data["json"] = theme
	c.ServeJson()
}

func (c *ThemeController) Put() {
	if c.user == nil || !c.user.IsManager() {
		jsonErr := simplejson.New()
		jsonErr.Set("error", "forbidden")
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(403)
		c.ServeJson()
		return
	}
	theme := &models.Theme{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), theme); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	theme, err := db.EditAndGetTheme(theme)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	c.Data["json"] = theme
	c.ServeJson()
}

func (c *ThemeController) Delete() {
	if c.user == nil || !c.user.IsManager() {
		jsonErr := simplejson.New()
		jsonErr.Set("error", "forbidden")
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(403)
		c.ServeJson()
		return
	}
	themeID, err := c.GetInt(":id", -1)
	if err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	if err := db.DeleteThemeByID(themeID); err != nil {
		beego.Error(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.Ctx.Output.SetStatus(400)
		c.ServeJson()
		return
	}
	c.Data["json"] = jsonOK
	c.ServeJson()
}
