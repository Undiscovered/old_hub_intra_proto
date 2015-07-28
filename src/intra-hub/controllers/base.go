package controllers

import (
	"encoding/gob"
	"strings"

	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const (
	sessionUserKey = "0xd3ob4"
)

func init() {
	gob.Register(&models.User{})
}

var (
	jsonOK = map[string]interface{}{"status": "OK"}
)

type NestedPreparer interface {
	NestedPrepare()
}

type BaseController struct {
	beego.Controller

	currentLanguage string
	user            *models.User
	isLogged        bool
	apiMode         bool
	flash           *beego.FlashData
	redirectURL     string
}

func (c *BaseController) Prepare() {
	// Set Default redirect URL
	c.redirectURL = c.Ctx.Input.Request.URL.Path

	// Set Language
	c.currentLanguage = "fr-FR"
	c.Data["Lang"] = c.currentLanguage

	// Set Flash data
	c.flash = beego.ReadFromRequest(&c.Controller)

	// Set the API mode if necessary.
	if strings.Contains(c.Ctx.Request.URL.Path, "/api") {
		c.apiMode = true
		c.EnableRender = false
	}

	// Get the user from the session.
	if user := c.GetSession(sessionUserKey); user != nil {
		c.user = user.(*models.User)
		c.isLogged = true
		c.Data["UserLogged"] = c.user
	}

	// If the matching controller is a NestedPreparer, we call the NestedPrepare function
	// To ensure that this Prepare function is called first (it is done to prevent overriding of
	// Prepare functions.
	if app, ok := c.AppController.(NestedPreparer); ok {
		app.NestedPrepare()
	}
}

func (c *BaseController) SetUser(user *models.User) {
	c.SetSession(sessionUserKey, user)
}

func (c *BaseController) SetErrorAndRedirect(err error) {
	c.flash.Data["error"] = err.Error()
	c.flash.Store(&c.Controller)
	c.Redirect(c.redirectURL, 303)
}

func (c *BaseController) RequireLogin() {
	if !c.isLogged || c.user == nil {
		c.Redirect("/login", 301)
	}
}

func (c *BaseController) RequireAdmin() {
	c.RequireLogin()
	if !c.user.IsAdmin() {
		c.flash.Data["error"] = "forbidden"
		c.flash.Store(&c.Controller)
		c.Redirect("/", 303)
	}
}

func (c *BaseController) RequireManager() {
	c.RequireLogin()
	if !c.user.IsManager() {
		c.flash.Data["error"] = "forbidden"
		c.flash.Store(&c.Controller)
		c.Redirect("/", 303)
	}
}

func (c *BaseController) RequirePedago() {
	c.RequireLogin()
	if !c.user.IsPedago() {
		c.flash.Data["error"] = "forbidden"
		c.flash.Store(&c.Controller)
		c.Redirect("/", 303)
	}
}
func (c *BaseController) TranslateSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = i18n.Tr(c.currentLanguage, s)
	}
	return slice
}
