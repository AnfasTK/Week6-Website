package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/routes"
)

var r *gin.Engine

func init() {
	r = gin.Default()
	r.LoadHTMLGlob("views/*")
	r.Static("/assets", "./assets")
	initializers.LoadEnvVariables()
	initializers.DBconnect()
	initializers.SyncDatabase()
}

func main() {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	routes.Routes(r)
	r.Run()
}
