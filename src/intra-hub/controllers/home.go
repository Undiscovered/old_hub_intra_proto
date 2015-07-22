package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) RootRedirection() {
	if c.isLogged {
		c.Redirect("/home", 301)
		return
	} else {
		c.Redirect("/login", 301)
		return
	}
}

func (c *HomeController) HomeView() {
	c.RequireLogin()
	c.TplNames = "list.html"
	c.Data["User"] = c.user
}
