package controllers

type UserController struct {
    BaseController
}

func (c *UserController) Login() {

}

func (c *UserController) LoginView() {
    c.TplNames = "login.html"
}
