package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

type HandlerI interface {
	CreateShortUrl(c *gin.Context)
	HandleShortUrlRedirect(c *gin.Context)
}

type handler struct {
	shortener shortener.ShortenerI
	cfg       *config.Config
	store     store.StorageServiceI
}

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func NewHandler(shortener shortener.ShortenerI, cfg *config.Config, store store.StorageServiceI) HandlerI {
	return &handler{
		shortener: shortener,
		cfg:       cfg,
		store:     store,
	}
}

func (h *handler) CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		return
	}
	h.store.SaveUrlMapping(c, shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	host := fmt.Sprintf("http://%s:%s/", h.cfg.ServerHost, h.cfg.ServerPort)
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func (h *handler) HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := h.store.RetrieveInitialUrl(c, shortUrl)
	c.Redirect(302, initialUrl)
}
