package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"main.go/auth"
	"main.go/initializers"
	"main.go/models"
)

var Error string
var FetchUser models.UserModel

func UserLoginHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	admincheck := session.Get(RoleAdmin)
	usercheck := session.Get(RoleUser)
	if usercheck !=nil && admincheck != nil {
		c.Redirect(http.StatusSeeOther, "/home")
	}else if admincheck == nil && usercheck !=nil   {
		c.Redirect(http.StatusSeeOther, "/home")
	} else if admincheck != nil && usercheck ==nil{
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	}else{
		c.HTML(http.StatusSeeOther, "userlogin.html", Error)
		Error = ""
	}
}

func UserLoginPostHandler(c *gin.Context) {

	FetchUser = models.UserModel{}
	err := initializers.DB.First(&FetchUser, "email=?", c.Request.FormValue("email")).Error
	if err != nil {
		Error = "invalid email Address"
		c.Redirect(http.StatusSeeOther, "/user/login")
	} else {
		plainpassword := c.Request.FormValue("password")

		err := bcrypt.CompareHashAndPassword([]byte(FetchUser.Password), []byte(plainpassword))

		if err != nil {
			Error = "invalid password"
			c.Redirect(http.StatusSeeOther, "/user/login")
		} else {
			if FetchUser.Status == "Blocked" {
				Error = "Your account has been blocked."
				c.Redirect(http.StatusSeeOther, "/user/login")
			} else {
				auth.JwtTokens(c, FetchUser.Email, RoleUser)
				c.Redirect(http.StatusSeeOther, "/home")
			}

		}
	}

}

func UserSignupHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	admincheck := session.Get(RoleAdmin)
	usercheck := session.Get(RoleUser)
	if usercheck !=nil && admincheck != nil {
		c.Redirect(http.StatusSeeOther, "/home")
	}else if admincheck == nil && usercheck !=nil   {
		c.Redirect(http.StatusSeeOther, "/home")
	} else if admincheck != nil && usercheck ==nil{
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	}else{
		c.HTML(http.StatusSeeOther, "usersignup.html", Error)
		Error = ""
	}
}

func UserSignupPostHandler(c *gin.Context) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Request.PostFormValue("password")), 10)

	if err != nil {
		log.Fatal("Error in password Hashing", err)
	}
	error := initializers.DB.Create(&models.UserModel{
		Name:     c.Request.PostFormValue("username"),
		Email:    c.Request.PostFormValue("email"),
		Password: string(hashedPassword),
		Status:   "Active",
	})
	if error.Error != nil {
		Error = "Email already exists"
		c.Redirect(http.StatusSeeOther, "/user/signup")
	} else {
		Error = "User Successfully Created"
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
}

func UserHomeHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleUser)
	if check != nil {
		c.HTML(http.StatusSeeOther, "userhome.html", FetchUser)
	} else {
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
}

func UserLogoutHandler(c *gin.Context) {
	c.Header("cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	session.Delete(RoleUser)
	session.Save()
	FetchUser = models.UserModel{}
	c.Redirect(http.StatusSeeOther, "/user/login")
}