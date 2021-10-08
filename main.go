package main

import (
	"github.com/geforce6t/go-server/models"
	"github.com/geforce6t/go-server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// initialise routes
	routes.InitRoutes(r)

	// initialise database
	models.InitDB()

	r.Use(cors.Default())
	r.Run(":6000")
}
