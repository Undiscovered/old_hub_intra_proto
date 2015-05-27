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
    if !c.isLogged {
        c.Redirect("/login", 301)
    }
    c.TplNames = "index.html"
    c.Data["User"] = c.user
}