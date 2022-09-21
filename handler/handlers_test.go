package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateShortUrl(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../dev.application.yml")
	if err != nil {
		return
	}
	ctx := context.Background()
	store := store.NewStorageService(cfg, ctx)
	h := handler.NewHandler(shortener, cfg, store)
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
