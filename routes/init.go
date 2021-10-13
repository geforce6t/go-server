package routes

import (
	"github.com/geforce6t/go-server/controllers"
	"github.com/geforce6t/go-server/middlewares"
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

	route.POST("/hello", middlewares.TokenAuthMiddleware(), func(c *gin.Context) {
		controllers.SayHello(c, db)
	})

	route.POST("/refresh", middlewares.TokenAuthMiddleware(), func(c *gin.Context) {
		controllers.Refresh(c, db)
	})
}
