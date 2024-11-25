package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var SecretKey []byte

type User struct {
	Userinfo string `json:"userinfo"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func InitSecretKey() {
	secret := os.Getenv("SECRET_KEY")
	if len(secret) == 0 {
		log.Fatal("SECRET_KEY environment variable is not set")
	}
	SecretKey = []byte(secret)
}

func JwtTokens(c *gin.Context, userinfo string, role string) {
	tokenString, err := CreateToken(userinfo, role)
	if err != nil {
		fmt.Println("Failed to create newtoken")
	}
	session := sessions.Default(c)
	session.Set(role, tokenString)
	session.Save()
	check := session.Get(role)
	fmt.Println("Jwt Token : ", check)
}

func CreateToken(userinfo string, role string) (string, error) {
	claims := User{
		Userinfo: userinfo,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", nil
	} else {
		return tokenString, nil
	}
}
