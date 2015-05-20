package routers

import (
	"intra-hub/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/api/login", &controllers.UserController{}, "post:Login")
    beego.Router("/login", &controllers.UserController{}, "get:LoginView;post:Login")
}
