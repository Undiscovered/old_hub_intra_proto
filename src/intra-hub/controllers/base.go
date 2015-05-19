package controllers

import (
	"github.com/astaxie/beego"
	"intra-hub/models"
)

const (
    sessionUserKey = "0xd3ob4"
)
type NestedPreparer interface{
    NestedPrepare()
}

type BaseRouter struct {
	beego.Controller

	user    *models.User
	isLogin bool
}

func (c *BaseRouter) Prepare() {
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

func (c *BaseRouter) SetUser(user *models.User) {
    c.SetSession(sessionUserKey, user)
}