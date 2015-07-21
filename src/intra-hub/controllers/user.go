package controllers

import (
	"fmt"
	"intra-hub/db"
	"intra-hub/models"
	"intra-hub/services/mail"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/bitly/go-simplejson"
	"github.com/jmcvetta/randutil"
	"intra-hub/jsonutils"
	"strconv"
	"strings"
)

type UserController struct {
	BaseController
}

func (c *UserController) AddView() {
	c.RequireAdmin()
	c.TplNames = "admin/add-user.html"
	groups, err := db.GetEveryGroups()
	if err != nil {
		c.Redirect("/", 301)
		return
	}
	cities, err := db.GetEveryCities()
	if err != nil {
		c.Redirect("/", 301)
		return
	}
	promotions, err := db.GetEveryPromotion()
	if err != nil {
		c.Redirect("/", 301)
		return
	}
	c.Data["Cities"] = cities
	c.Data["Groups"] = groups
	c.Data["Promotions"] = promotions
}

func (c *UserController) SingleView() {
	c.RequireLogin()
	c.TplNames = "user/profile.html"
	user, err := db.GetUserByLogin(c.GetString(":id", ""))
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
		return
	}
	c.Data["User"] = user
	c.Data["UserJSON"] = user.ToJSON(c.currentLanguage)
}

func (c *UserController) MeView() {
	c.RequireLogin()
	c.TplNames = "user/profile.html"
	c.Data["User"] = c.user
	c.Data["UserJSON"] = c.user.ToJSON(c.currentLanguage)
	c.Data["Edit"] = true
}

func (c *UserController) LoginView() {
	c.TplNames = "login.html"
}

func (c *UserController) EditView() {
	c.RequireLogin()
	c.TplNames = "user/edit.html"
	user, err := db.GetUserByLogin(c.GetString(":login", ""))
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	if c.user.Login != user.Login && !c.user.IsManager() {
		c.Redirect("/home", 301)
	}
	cities, err := db.GetEveryCities()
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	skills, err := db.GetEverySkills()
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	groups, err := db.GetEveryGroups()
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	themes, err := db.GetEveryThemes()
	if err != nil {
		beego.Error(err)
		c.Redirect("/home", 301)
	}
	c.Data["User"] = user
	c.Data["Edit"] = user.Login == c.user.Login
	c.Data["Cities"] = jsonutils.MarshalUnsafe(cities)
	c.Data["Groups"] = jsonutils.MarshalUnsafe(groups)
	c.Data["UserJSON"] = user.ToJSON(c.currentLanguage)
	c.Data["Skills"] = jsonutils.MarshalUnsafe(skills)
	c.Data["Themes"] = jsonutils.MarshalUnsafe(themes)
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
	c.SetUser(user)
	c.Redirect("/home", 301)
}

func (c *UserController) Logout() {
	c.DestroySession()
	c.Redirect("/login", 301)
}

func (c *UserController) SearchUser() {
	c.RequireLogin()
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
	token := c.GetString(":token")
	id, err := strconv.Atoi(c.GetString(":id"))
	if err != nil {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
		return
	}
	if err := db.CheckUserExists(id, token); err != nil {
		beego.Error(err)
		c.DestroySession()
		c.Redirect("/login", 301)
		return
	}
	c.Data["Id"] = id
	c.Data["Token"] = token
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
	c.RequireManager()
	c.redirectURL = "/admin/users/add"
	user := &models.User{}
	if err := c.ParseForm(user); err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	if user.FirstName == "" {
		c.SetErrorAndRedirect(fmt.Errorf("first name not specified"))
		return
	}
	if user.LastName == "" {
		c.SetErrorAndRedirect(fmt.Errorf("last name not specified"))
		return
	}
	if user.Login == "" {
		c.SetErrorAndRedirect(fmt.Errorf("login not specified"))
		return
	}
	if user.Email == "" {
		c.SetErrorAndRedirect(fmt.Errorf("wrong email format"))
		return
	}
	user.Promotion = &models.Promotion{Id: user.PromotionID}
	user.Group = &models.Group{Id: user.GroupID}
	user.City = &models.City{Id: user.CityID}
	randString, err := randutil.AlphaString(9)
	if err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	user.Token = randString
	if err := db.AddUser(user); err != nil {
		c.SetErrorAndRedirect(err)
		return
	}
	go mail.SendUserCreated(user)
	c.flash.Data["success"] = "User created"
	c.Redirect("/admin/users/add", 301)
}

func (c *UserController) EditUser() {
	c.RequireLogin()
	c.EnableRender = false
	user := &models.User{}
	if err := json.Unmarshal(c.Ctx.Input.CopyBody(), &user); err != nil {
		beego.Warn(err)
		return
	}
	if c.user.Login != user.Login && !c.user.IsManager() {
		c.Redirect("/home", 301)
	}
	userDB, err := db.GetUserByLogin(user.Login)
	if err != nil {
		beego.Warn(err)
		return
	}
	if !c.user.IsAdmin() {
		user.Group = userDB.Group
	}
	if err := db.EditUserByLogin(user.Login, user); err != nil {
		beego.Warn(err)
		return
	}
	if user.Login == c.user.Login {
		c.SetUser(user)
	}
}

func (c *UserController) GetMe() {
	defer c.ServeJson()
	c.Data["json"] = c.user.Clean()
}

func (c *UserController) ListStudentView() {
	c.TplNames = "student/list.html"
	handleError := func(err error) {
		beego.Error(err)
		c.SetErrorAndRedirect(err)
	}
	promotions, err := db.GetEveryPromotion()
	if err != nil {
		handleError(err)
		return
	}
	cities, err := db.GetEveryCities()
	if err != nil {
		handleError(err)
		return
	}
	skills, err := db.GetEverySkills()
	if err != nil {
		handleError(err)
		return
	}
	themes, err := db.GetEveryThemes()
	if err != nil {
		handleError(err)
		return
	}
	queryFilter := make(map[string]interface{})
	queryFilter["promotions"] = strings.Split(c.GetString("promotions", ""), ",")
	queryFilter["cities"] = strings.Split(c.GetString("cities", ""), ",")
	queryFilter["skills"] = strings.Split(c.GetString("skills", ""), ",")
	queryFilter["themes"] = strings.Split(c.GetString("themes", ""), ",")
	queryFilter["name"] = c.GetString("name", "")
	queryFilter["login"] = c.GetString("login", "")
	queryFilter["email"] = c.GetString("email", "")
	page, err := c.GetInt("page")
	if err != nil {
		handleError(err)
		return
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		handleError(err)
		return
	}
	if page <= 0 {
		c.Redirect(fmt.Sprintf("/students?page=1&limit=%d", limit), 301)
		return
	}
	if limit == 0 {
		limit = 25
	}
	paginatedItems, err := db.GetUsersPaginated(page, limit, queryFilter)
	if err != nil {
		handleError(err)
		return
	}
	paginatedItems.SetPagesToShow()
	c.Data["Cities"] = cities
	c.Data["Promotions"] = promotions
	c.Data["Skills"] = skills
	c.Data["Themes"] = themes
	c.Data["Limit"] = limit
	c.Data["PaginatedItems"] = paginatedItems
	c.Data["HasNextPage"] = paginatedItems.CurrentPage+1 <= paginatedItems.TotalPageCount
	c.Data["HasPreviousPage"] = paginatedItems.CurrentPage != 1
	c.Data["ShowGoToFirst"] = paginatedItems.PagesToShow[0] != 1
	c.Data["ShowGoToLast"] = paginatedItems.PagesToShow[len(paginatedItems.PagesToShow)-1] != paginatedItems.TotalPageCount
}
