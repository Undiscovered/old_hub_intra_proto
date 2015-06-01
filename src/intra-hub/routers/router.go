package routers

import (
	"intra-hub/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.HomeController{}, "get:RootRedirection")
    beego.Router("/home", &controllers.HomeController{}, "get:HomeView")
    beego.Router("/logout", &controllers.UserController{}, "get:Logout")
    beego.Router("/login", &controllers.UserController{}, "get:LoginView;post:Login")
    beego.Router("/users/search", &controllers.UserController{}, "post:SearchUser")
    beego.Router("/projects", &controllers.ProjectController{}, "get:IntroView")
    beego.Router("/projects/list", &controllers.ProjectController{}, "get:ListView")
    beego.Router("/projects/:nameOrId", &controllers.ProjectController{}, "get:SingleView")
    beego.Router("/projects/add", &controllers.ProjectController{}, "get:AddView;post:Add")

    beego.Router("/api/login", &controllers.UserController{}, "post:Login")
}
