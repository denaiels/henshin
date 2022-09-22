package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/handler"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func MockJSONPost(c *gin.Context, urlCreationRequest handler.UrlCreationRequest) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(urlCreationRequest)
	if err != nil {
		log.Err(err).Msg("Error while mocking JSON post")
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateShortUrl(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
