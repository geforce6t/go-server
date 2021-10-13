package main

import (
	"github.com/geforce6t/go-server/models"
	"github.com/geforce6t/go-server/routes"
	"github.com/geforce6t/go-server/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// initialise database
	db := models.InitDB()

	// initialise redis
	utils.InitialiseRedis()

	// initialise routes
	routes.InitRoutes(r, db)

	r.Use(cors.Default())
	r.Run(":6000")
}
