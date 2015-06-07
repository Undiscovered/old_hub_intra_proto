package controllers
import (
    "intra-hub/models"
    "github.com/astaxie/beego"
    "intra-hub/db"
    "encoding/json"
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
        jsonErr.Set("error", err)
        c.Data["json"] = jsonErr
        c.ServeJson()
        return
    }
    c.Data["json"] = theme
    c.ServeJson()
}