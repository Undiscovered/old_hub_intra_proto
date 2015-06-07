package controllers
import (
    "intra-hub/db"
    "github.com/astaxie/beego"
)

type AdminController struct {
    BaseController
}

func (c *AdminController) Get() {
    c.TplNames = "admin/layout.html"
    themes, err := db.GetAllThemes()
    if err != nil {
        beego.Error(err)
        c.SetErrorAndRedirect(err)
        return
    }
    beego.Warn(themes)
    c.Data["Themes"] = themes
}