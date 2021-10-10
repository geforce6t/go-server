package utils

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geforce6t/go-server/models"
)

func GenerateJwt(id uint, name string) (string, error) {
	var mySigningKey = []byte(models.GetEnvValue("secret"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatalf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
