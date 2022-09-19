package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"source.golabs.io/daniel.santoso/url-blaster/handler"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func main() {
	shortener := shortener.NewShortener()
	handler := handler.NewHandler(shortener)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is the Go URL Blaster!",
		})
	})

	router.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	router.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	store.InitializeStore()

	err := StartWebServer(router, "9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error %v", err))
	}
}

func StartWebServer(router *gin.Engine, portNumber string) error {
	err := router.Run(fmt.Sprintf(":%s", portNumber))
	return err
}
