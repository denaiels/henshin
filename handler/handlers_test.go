package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func MockJSONPost(c *gin.Context, urlCreationRequest UrlCreationRequest) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(urlCreationRequest)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreatingShortUrl(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJSONPost(c, UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	store.InitializeStore()
	CreateShortUrl(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
