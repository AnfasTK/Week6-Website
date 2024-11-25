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

var AdError string
var FetchAdmin models.AdminModel
var UserData []models.UserModel
var UserUpdate models.UserModel

func AdminLoginHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	usercheck := session.Get(RoleUser)
	admincheck := session.Get(RoleAdmin)
	if usercheck != nil && admincheck != nil {
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	} else if admincheck != nil && usercheck == nil {
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	} else if admincheck == nil && usercheck != nil {
		c.Redirect(http.StatusSeeOther, "/home")
	} else {
		c.HTML(http.StatusSeeOther, "adminlogin.html", AdError)
		AdError = ""
	}

}

func AdminLoginPostHandler(c *gin.Context) {
	FetchAdmin = models.AdminModel{}
	err := initializers.DB.First(&FetchAdmin, "email=?", c.Request.PostFormValue("email")).Error
	if err != nil {
		AdError = "Invalid Email Address"
		c.Redirect(http.StatusSeeOther, "/admin/login")
	} else {
		if FetchAdmin.Password != c.Request.FormValue("password") {
			AdError = "Invalid password"
			c.Redirect(http.StatusSeeOther, "/admin/login")
		} else {
			auth.JwtTokens(c, FetchAdmin.Email, RoleAdmin)
			c.Redirect(http.StatusSeeOther, "/adminpannel")
		}
	}

}

func AdminPannelHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		initializers.DB.Find(&UserData)
		c.HTML(http.StatusSeeOther, "adminhome.html", gin.H{
			"UserDatas": UserData,
			"Admin":     FetchAdmin.Name,
			"Error":     AdError,
		})
		AdError = ""
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func AdminSearchHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)

	if check != nil {

		query := c.Query("query")

		if query == "" {
			AdError = "Search cannot be empty."
			c.Redirect(http.StatusSeeOther, "/adminpannel")
			return
		}

		var matchedUsers []models.UserModel
		 searchQuery := query + "%"
		 err := initializers.DB.Where("name iLIKE ? OR email iLIKE ?", searchQuery, searchQuery).Find(&matchedUsers).Error

		if err != nil {
			AdError = "Error fetching search results."
			c.Redirect(http.StatusSeeOther, "/adminpannel")
			return
		}

		c.HTML(http.StatusOK, "adminhome.html", gin.H{
			"UserDatas": matchedUsers,
			"Admin":     FetchAdmin.Name,
			"Error":     AdError,
		})
		AdError = ""
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func UserCreateHandler(c *gin.Context) {
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
		AdError = "User Not Created Email Already Existing"
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	} else {
		AdError = "User Successfully Created"
		c.Redirect(http.StatusSeeOther, "/adminpannel")
	}
}

func EditHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		UserId := c.Param("ID")
		c.HTML(http.StatusSeeOther, "edit.html", UserId)
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func EditPostHandler(c *gin.Context) {
	UserId := c.Param("ID")
	initializers.DB.Find(&UserUpdate, "ID=?", UserId)
	UserUpdate.Name = c.Request.FormValue("username")
	UserUpdate.Email = c.Request.FormValue("email")
	err := initializers.DB.Save(&UserUpdate).Error
	if err != nil {
		AdError = "Update Not Saved Email Already Existing"
	} else {
		AdError = "Update Saved"
	}
	UserUpdate = models.UserModel{}
	c.Redirect(http.StatusSeeOther, "/adminpannel")
}

func UserSatusHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		User := c.Param("ID")
		initializers.DB.First(&UserUpdate, "ID=?", User)

		if UserUpdate.Status == "Active" {
			UserUpdate.Status = "Blocked"
			initializers.DB.Save(&UserUpdate)
			UserUpdate = models.UserModel{}
			c.Redirect(http.StatusSeeOther, "/adminpannel")
			AdError = "You Have Successfully Blocked The User."
		} else {
			UserUpdate.Status = "Active"
			initializers.DB.Save(&UserUpdate)
			UserUpdate = models.UserModel{}
			c.Redirect(http.StatusSeeOther, "/adminpannel")
			AdError = "You Successfully Ativated The User."
		}
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func UserDeleteHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		User := c.Param("ID")
		initializers.DB.First(&UserUpdate, "ID=?", User)
		initializers.DB.Delete(&UserUpdate)
		UserUpdate = models.UserModel{}
		c.Redirect(http.StatusSeeOther, "/adminpannel")
		AdError = "You Have Successfully Deleted The User."
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/login")
	}

}

func AdminLogoutHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	session.Delete(RoleAdmin)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin/login")
}