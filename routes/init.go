package routes

import (
	"github.com/geforce6t/go-server/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(route *gin.Engine) {
	route.GET("hello", controllers.SayHello)
}
