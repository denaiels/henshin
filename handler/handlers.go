package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

type HandlerI interface {
	CreateShortUrl(c *gin.Context)
	HandleShortUrlRedirect(c *gin.Context)
}

type handler struct {
	cfg       *config.Config
	shortener shortener.ShortenerI
	store     store.StorageServiceI
}

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id"`
}

func NewHandler(shortener shortener.ShortenerI, cfg *config.Config, store store.StorageServiceI) HandlerI {
	return &handler{
		cfg:       cfg,
		shortener: shortener,
		store:     store,
	}
}

func (h *handler) CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !strings.HasPrefix(creationRequest.LongUrl, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please input a valid url!"})
		return
	}

	if creationRequest.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please input a valid user id!"})
		return
	}

	shortUrl, err := h.shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		log.Err(err).Msg("Error while generating short link because error while encoding with base58")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.store.SaveUrlMapping(c, shortUrl, creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s", err, shortUrl, creationRequest.LongUrl))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	host := fmt.Sprintf("http://%s:%s/", h.cfg.ServerHost, h.cfg.ServerPort)
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func (h *handler) HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl, err := h.store.RetrieveInitialUrl(c, shortUrl)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed retrieving inital url | Error: %v - shortUrl: %s", err, shortUrl))
		c.JSON(404, gin.H{
			"message": "Something's wrong, i can feel it... Maybe you entered the wrong link.",
		})
		return
	}
	c.Redirect(302, initialUrl)
}
