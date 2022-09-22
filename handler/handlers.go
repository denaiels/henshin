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
	UpdateLongUrl(c *gin.Context)
	HandleShortUrlRedirect(c *gin.Context)
	RemoveShortUrl(c *gin.Context)
}

type handler struct {
	cfg       *config.Config
	shortener shortener.ShortenerI
	store     store.StorageServiceI
}

type UrlCreationRequest struct {
	LongUrl        string `json:"long_url" binding:"required"`
	UserId         string `json:"user_id"`
	PredefinedName string `json:"predefined_name"`
}

type UrlUpdateRequest struct {
	ShortUrl   string `json:"short_url" binding:"required"`
	NewLongUrl string `json:"new_long_url" binding:"required"`
}

type UrlRemoveRequest struct {
	ShortUrl string `json:"short_url" binding:"required"`
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

	var shortUrl string
	var err error
	if creationRequest.PredefinedName != "" {
		shortUrl = creationRequest.PredefinedName
	} else {
		shortUrl, err = h.shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
		if err != nil {
			log.Err(err).Msg("Error while generating short link because error while encoding with base58")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err = h.store.SaveUrlMapping(c, shortUrl, creationRequest.LongUrl)
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

func (h *handler) UpdateLongUrl(c *gin.Context) {
	var updateRequest UrlUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !strings.HasPrefix(updateRequest.NewLongUrl, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please input a valid url!"})
		return
	}

	if !h.store.CheckIfShortUrlExists(c, updateRequest.ShortUrl) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short url doesn't exist!"})
		return
	}

	err := h.store.SaveUrlMapping(c, updateRequest.ShortUrl, updateRequest.NewLongUrl)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s", err, updateRequest.ShortUrl, updateRequest.NewLongUrl))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "url updated successfully",
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

func (h *handler) RemoveShortUrl(c *gin.Context) {
	var removeRequest UrlRemoveRequest
	if err := c.ShouldBindJSON(&removeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !h.store.CheckIfShortUrlExists(c, removeRequest.ShortUrl) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short url doesn't exist!"})
		return
	}

	err := h.store.DeleteUrlMapping(c, removeRequest.ShortUrl)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Failed deleting key url | Error: %v - shortUrl: %s", err, removeRequest.ShortUrl))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "short url deleted successfully",
	})
}
