package routers

import (
	"intra-hub/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "get:RootRedirection")

	beego.Router("/me", &controllers.UserController{}, "get:MeView")
	beego.Router("/home", &controllers.HomeController{}, "get:HomeView")
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")
	beego.Router("/login", &controllers.UserController{}, "get:LoginView;post:Login")
	beego.Router("/forgot", &controllers.UserController{}, "get:ResetPasswordView;post:ResetPassword")

	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/users/add", &controllers.UserController{}, "get:AddView")

	beego.Router("/theme", &controllers.ThemeController{})
	beego.Router("/theme/:id", &controllers.ThemeController{})

	beego.Router("/skills", &controllers.SkillController{})
	beego.Router("/skills/:id", &controllers.SkillController{})

	beego.Router("/staff", &controllers.StaffController{}, "get:ListView")

	beego.Router("/students", &controllers.UserController{}, "get:ListStudentView")

	beego.Router("/users/activate/:id/:token", &controllers.UserController{}, "get:ActivateUserView;post:ActivateUser")
	beego.Router("/users/edit/:login", &controllers.UserController{}, "get:EditView;put:EditUser")
	beego.Router("/users/search", &controllers.UserController{}, "post:SearchUser")
	beego.Router("/users/:id", &controllers.UserController{}, "get:SingleView")
	beego.Router("/users", &controllers.UserController{}, "post:AddUser")

	beego.Router("/projects", &controllers.ProjectController{}, "get:ListView")
	beego.Router("/projects/checkname", &controllers.ProjectController{}, "get:CheckName")
	beego.Router("/projects/add", &controllers.ProjectController{}, "get:AddView;post:Add")
	beego.Router("/projects/edit", &controllers.ProjectController{}, "post:Edit")
	beego.Router("/projects/comments", &controllers.ProjectController{}, "post:AddComment")
	beego.Router("/projects/edit/:nameOrId", &controllers.ProjectController{}, "get:EditView")
	beego.Router("/projects/:nameOrId", &controllers.ProjectController{}, "get:SingleView")
	beego.Router("/projects/:nameOrId/comments", &controllers.ProjectController{}, "get:CommentView;post:AddComment")

	beego.Router("/pedago/validate/:userId/:projectId/:validation", &controllers.PedagoController{}, "post:ValidateProject")
	beego.Router("/pedago/validation/:validation", &controllers.PedagoController{}, "get:ValidateProjectView")

	beego.Router("/export/projects", &controllers.ExportController{}, "get:Projects")

	beego.Router("/api/login", &controllers.UserController{}, "post:Login")
	beego.Router("/api/users/me", &controllers.UserController{}, "get:GetMe")
}
