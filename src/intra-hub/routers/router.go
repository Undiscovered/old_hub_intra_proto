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

	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/users/add", &controllers.UserController{}, "get:AddView")

	beego.Router("/theme", &controllers.ThemeController{})
	beego.Router("/theme/:id", &controllers.ThemeController{})

	beego.Router("/technos", &controllers.TechnoController{})
	beego.Router("/technos/:id", &controllers.TechnoController{})

	beego.Router("/skills", &controllers.SkillController{})
	beego.Router("/skills/:id", &controllers.SkillController{})

	beego.Router("/users/search", &controllers.UserController{}, "post:SearchUser")
	beego.Router("/users/:id", &controllers.UserController{}, "get:SingleView")

	beego.Router("/projects", &controllers.ProjectController{}, "get:IntroView")
	beego.Router("/projects/list", &controllers.ProjectController{}, "get:ListView")
	beego.Router("/projects/:nameOrId", &controllers.ProjectController{}, "get:SingleView")
	beego.Router("/projects/add", &controllers.ProjectController{}, "get:AddView;post:Add")

	beego.Router("/api/login", &controllers.UserController{}, "post:Login")
}
