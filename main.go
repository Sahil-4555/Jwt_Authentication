package main

import (
	"os"

	middleware "github.com/sahil/jwt-auth-go/middleware"
	routes "github.com/sahil/jwt-auth-go/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/sahil/jwt-auth-go/configs"
)

func main() {
	port := os.Getenv(configs.PORT())

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	// API-2
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
