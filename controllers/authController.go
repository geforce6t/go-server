package controllers

import "github.com/gin-gonic/gin"

type RegisterResponse struct {
	Name  string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
}

func RegisterUser(c *gin.Context) {
	c.ShouldBindJSON()

}

func LoginUser() {

}
