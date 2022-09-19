package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"source.golabs.io/daniel.santoso/url-blaster/handler"
	"source.golabs.io/daniel.santoso/url-blaster/shortener"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

type handlerSuite struct {
	suite.Suite
	urlCreationRequest *handler.UrlCreationRequest
}

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
	h := handler.NewHandler(shortener)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	store.InitializeStore()
	h.CreateShortUrl(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
