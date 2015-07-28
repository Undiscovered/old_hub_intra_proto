package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) HomeView() {
	c.RequireLogin()
	c.TplNames = "index.html"
	c.Data["User"] = c.user
}
