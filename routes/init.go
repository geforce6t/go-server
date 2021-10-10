package routes

import (
	"github.com/geforce6t/go-server/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(route *gin.Engine, db *gorm.DB) {
	route.POST("/register", func(c *gin.Context) {
		controllers.RegisterUser(c, db)
	})

	route.POST("/login", func(c *gin.Context) {
		controllers.LoginUser(c, db)
	})
}
