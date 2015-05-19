package routers

import (
	"intra-hub/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/login", &controllers.UserController{}, "get:LoginView;post:Login")
}
