package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "intra.hub.epitech.eu"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplNames = "index.html"
}
