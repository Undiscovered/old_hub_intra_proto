package controllers

import (
	"fmt"
	"strings"

	"intra-hub/db"
	"intra-hub/models"
	"intra-hub/services/mail"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/bitly/go-simplejson"
    "strconv"
    "github.com/jmcvetta/randutil"
)

type UserController struct {
	BaseController
}

func (c *UserController) AddView() {
	c.TplNames = "admin/add-user.html"
	c.Data["Groups"] = models.EveryUserGroups
}

func (c *UserController) SingleView() {
	c.TplNames = "user/profile.html"
	user, err := db.GetUserByLogin(c.GetString(":id", ""))
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
		return
	}
	c.Data["User"] = user
}

func (c *UserController) LoginView() {
	c.TplNames = "login.html"
}

func (c *UserController) EditView() {
	c.TplNames = "user/edit.html"
	user, err := db.GetUserByLogin(c.GetString(":login", ""))
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	skills, err := db.GetEverySkills()
	c.Data["Groups"] = models.EveryUserGroups
	c.Data["User"] = user
	c.Data["Skills"] = skills
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

func (c *UserController) SearchUser() {
	c.EnableRender = false
	jsonBody, err := simplejson.NewJson(c.Ctx.Input.CopyBody())
	if err != nil {
		beego.Warn(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
		c.ServeJson()
		return
	}
	users, err := db.SearchUsers(jsonBody.Get("login").MustString(""))
	if err != nil {
		beego.Warn(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
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

func (c *UserController) ActivateUserView() {
	c.TplNames = "user/reset-password.html"
    id, err := strconv.Atoi(c.GetString(":id"))
    if err != nil {
        beego.Error(err)
        c.SetErrorAndRedirect(err)
        return
    }
    c.Data["Id"] = id
    c.Data["Token"] = c.GetString(":token")
}

func (c *UserController) ActivateUser() {
    c.EnableRender = false
    password := c.GetString("password")
    token := c.GetString(":token")
    id, err := strconv.Atoi(c.GetString(":id"))
    if err != nil {
        beego.Error(err)
        c.SetErrorAndRedirect(err)
        return
    }
    user, err := db.ActivateUser(id, token, password)
    if err != nil {
        beego.Error(err)
        c.SetErrorAndRedirect(err)
        return
    }
    c.SetUser(user)
    c.Redirect("/home", 301)
}

func (c *UserController) AddUser() {
	defer c.ServeJson()
	handleError := func(err error) {
		beego.Warn(err)
		jsonErr := simplejson.New()
		jsonErr.Set("error", err.Error())
		c.Data["json"] = jsonErr
	}
	c.EnableRender = false
	user := &models.User{}
	if err := c.ParseForm(user); err != nil {
		handleError(err)
		return
	}
	if user.Login == "" {
		handleError(fmt.Errorf("login not specified"))
		return
	}
	if user.Email == "" || !strings.Contains(user.Email, "@") {
		handleError(fmt.Errorf("wrong email format"))
		return
	}
    randString, err := randutil.AlphaString(9)
    if err != nil {
        handleError(err)
        return
    }
    user.Token = randString
	if err := db.AddUser(user); err != nil {
		handleError(err)
		return
	}
	go mail.SendUserCreated(user)
}
