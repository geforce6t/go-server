package controllers

import "github.com/gin-gonic/gin"

func SayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "posted",
		"message": "hello",
	})
}
