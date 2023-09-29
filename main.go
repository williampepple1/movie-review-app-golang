package main

import (
	"os"
	"review-app/database"
	"review-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.Default()

	//run database
	database.StartDB()

	//Log events
	router.Use(gin.Logger())

	//Register app routes here
	routes.AuthRoutes(router)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Welcome to shive api!"})
	})

	router.Run(":" + port)

}