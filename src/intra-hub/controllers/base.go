package controllers

import (
	"github.com/astaxie/beego"
	"intra-hub/models"
    "strings"
)

const (
	sessionUserKey = "0xd3ob4"
)

type NestedPreparer interface {
	NestedPrepare()
}

type BaseController struct {
	beego.Controller

	user    *models.User
	isLogin bool
	apiMode bool
}

func (c *BaseController) Prepare() {
    // Set the API mode if necessary.
    if strings.Contains(c.Ctx.Request.URL.Path, "/api") {
        c.apiMode = true
        c.EnableRender = false
    }

	// Get the user from the session.
	if user := c.GetSession(sessionUserKey); user != nil {
		c.user = user.(*models.User)
		c.isLogin = true
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
