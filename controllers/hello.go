package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SayHello(c *gin.Context, db *gorm.DB) {
	c.JSON(http.StatusAccepted, gin.H{
		"message": "contratulations! Here we go",
	})
}
