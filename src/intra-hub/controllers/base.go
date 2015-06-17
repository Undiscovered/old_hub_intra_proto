package controllers

import (
	"math/rand"
	"strconv"
	"strings"

	"intra-hub/models"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const (
	sessionUserKey = "0xd3ob4"
)

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

	// Add basic template functions
	beego.AddFuncMap("i18n", i18n.Tr)
	incr := func(arg int) string {
		return strconv.FormatInt(int64(arg+1), 10)
	}
	decr := func(arg int) string {
		return strconv.FormatInt(int64(arg-1), 10)
	}
	randomizeLabel := func() string {
		labels := []string{"success", "warning", "danger", "info", "primary", "default"}
		return labels[rand.Intn(len(labels))]
	}
	beego.AddFuncMap("incr", incr)
	beego.AddFuncMap("decr", decr)
	beego.AddFuncMap("randLabel", randomizeLabel)

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
    if !c.isLogged {
        c.Redirect("/login", 301)
    }
}

func (c *BaseController) TranslateSlice(slice []string) []string {
    for i, s := range slice {
        slice[i] = i18n.Tr(c.currentLanguage, s)
    }
    return slice
}