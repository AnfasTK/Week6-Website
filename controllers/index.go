package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const RoleAdmin = "admin"
const RoleUser = "user"

func IndexHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	session := sessions.Default(c)
	userCheck := session.Get(RoleUser)
	adminCheck := session.Get(RoleAdmin)
	if userCheck == nil && adminCheck != nil {
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	} else if adminCheck == nil && userCheck != nil {
		c.Redirect(http.StatusSeeOther, "/home")
	} else {
		c.HTML(http.StatusOK, "index.html", nil)
	}

}
