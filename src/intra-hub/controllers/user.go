package controllers

import (
	"github.com/astaxie/beego/validation"
	"github.com/bitly/go-simplejson"
	"intra-hub/db"
	"intra-hub/models"
    "github.com/astaxie/beego"
    "fmt"
)

type UserController struct {
	BaseController
}

func (c *UserController) Login() {
    user := &models.User{}
	if c.isLogged {
		c.Redirect("/home", 301)
		return
	}
	c.ParseForm(user)
	valid := validation.Validation{}
	if b, err := valid.Valid(user); err != nil {
        c.SetErrorAndRedirect(err)
		return
	} else if !b {
        c.SetErrorAndRedirect(fmt.Errorf(valid.Errors[0].Message))
		return
	}
	user, err := db.CheckUserCredentials(user)
	if err != nil {
        c.SetErrorAndRedirect(err)
		return
	}
    beego.Alert(user.Promotion, user.City)
	c.SetUser(user)
	c.Redirect("/home", 301)
}

func (c *UserController) Logout() {
    c.DestroySession()
    c.Redirect("/login", 301)
}

func (c *UserController) LoginView() {
	c.TplNames = "login.html"
}

func (c *UserController) SearchUser() {
	c.EnableRender = false
	jsonBody, err := simplejson.NewJson(c.Ctx.Input.CopyBody())
	if err != nil {
        beego.Warn(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	users, err := db.SearchUsers(jsonBody.Get("login").MustString(""))
	if err != nil {
        beego.Warn(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err)
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
    if users == nil {
        // Return an empty array.
        c.Data["json"] = []string{}
    } else {
        c.Data["json"] = users
    }
	c.ServeJson()
}

func (c *UserController) AddView() {

}