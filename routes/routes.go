package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controllers"
)

func Routes(r *gin.Engine) {
	// Index page
	r.GET("/", controllers.IndexHandler)

	//User
	r.GET("/user/login", controllers.UserLoginHandler)
	r.POST("/login", controllers.UserLoginPostHandler)
	r.GET("/user/signup", controllers.UserSignupHandler)
	r.POST("/signup", controllers.UserSignupPostHandler)
	r.GET("/home", controllers.UserHomeHandler)
	r.GET("/logout", controllers.UserLogoutHandler)

	//Admin
	r.GET("/admin/login", controllers.AdminLoginHandler)
	r.POST("/adminlogin", controllers.AdminLoginPostHandler)
	r.GET("/adminpannel", controllers.AdminPannelHandler)
	r.GET("/admin/search", controllers.AdminSearchHandler)
	r.POST("/admin/create-user", controllers.UserCreateHandler)
	r.GET("/admin/edit-user/:ID", controllers.EditHandler)
	r.POST("/admin/edit-user/:ID", controllers.EditPostHandler)
	r.GET("/admin/block-user/:ID", controllers.UserSatusHandler)
	r.GET("/admin/delete-user/:ID", controllers.UserDeleteHandler)
	r.GET("/adminlogout", controllers.AdminLogoutHandler)
}
