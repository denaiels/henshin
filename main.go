package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is the Go URL Blaster!",
		})
	})

	err := StartWebServer(router, "9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error %v", err))
	}
}

func StartWebServer(router *gin.Engine, portNumber string) error {
	err := router.Run(fmt.Sprintf(":%s", portNumber))
	return err
}
