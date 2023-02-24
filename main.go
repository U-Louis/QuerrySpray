package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define an endpoint that returns a JSON response
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	// Start the server
	if err := router.Run(":8085"); err != nil {
		panic("Failed to start server")
	}
}