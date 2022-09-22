package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/handler"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func main() {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("dev.application.yml")
	if err != nil {
		log.Err(err).Msg("Error while loading config")
	}
	ctx := context.Background()
	store := store.NewStorageService(cfg, ctx)
	handler := handler.NewHandler(shortener, cfg, store)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is the Go URL Blaster!",
		})
	})

	router.POST("/create-short-url", handler.CreateShortUrl)

	router.POST("/update-url", handler.UpdateLongUrl)

	router.POST("/remove-url", handler.RemoveShortUrl)

	router.GET("/:shortUrl", handler.HandleShortUrlRedirect)

	err = StartWebServer(router, cfg.ServerPort)
	if err != nil {
		log.Panic().Msg(fmt.Sprintf("Failed to start the web server - Error %v", err))
	}
}

func StartWebServer(router *gin.Engine, portNumber string) error {
	err := router.Run(fmt.Sprintf(":%s", portNumber))
	return err
}
