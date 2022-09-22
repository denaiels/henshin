package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"
const CacheDuration = 6 * time.Hour

func MockCreationJSONPost(c *gin.Context, urlCreationRequest handler.UrlCreationRequest) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(urlCreationRequest)
	if err != nil {
		log.Err(err).Msg("Error while mocking JSON post")
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func MockUpdateJSONPost(c *gin.Context, updateRequest handler.UrlUpdateRequest) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(updateRequest)
	if err != nil {
		log.Err(err).Msg("Error while mocking JSON post")
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func MockRemoveJSONPost(c *gin.Context, removeRequest handler.UrlRemoveRequest) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(removeRequest)
	if err != nil {
		log.Err(err).Msg("Error while mocking JSON post")
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateShortUrlSuccess(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateShortUrlWithPredefinedNameSuccess(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl:        "https://youtu.be/8LhMu4bQTQU",
		UserId:         "e0dba740-fc4b-4977-872c-d360239e6b10",
		PredefinedName: "dyna",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateShortUrlEmptyUrl(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShortUrlEmptyUserId(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShortUrlWithInvalidLink(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "hahaha",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShortUrlRedisFail(t *testing.T) {
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

	MockCreationJSONPost(c, handler.UrlCreationRequest{
		LongUrl: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		UserId:  "e0dba740-fc4b-4977-872c-d360239e6b10",
	})

	redisServer.SetError("REDISDOWN")
	h.CreateShortUrl(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateUrlSuccess(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockUpdateJSONPost(c, handler.UrlUpdateRequest{
		ShortUrl:   "dyna",
		NewLongUrl: "https://youtu.be/UIbNIhaldLQ",
	})

	h.UpdateLongUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUrlEmptyShortUrl(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockUpdateJSONPost(c, handler.UrlUpdateRequest{
		ShortUrl:   "",
		NewLongUrl: "https://youtu.be/UIbNIhaldLQ",
	})

	h.UpdateLongUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUrlInvalidLink(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockUpdateJSONPost(c, handler.UrlUpdateRequest{
		ShortUrl:   "dyna",
		NewLongUrl: "hahahaha",
	})

	h.UpdateLongUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUrlRedisFail(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockUpdateJSONPost(c, handler.UrlUpdateRequest{
		ShortUrl:   "dyna",
		NewLongUrl: "https://youtu.be/UIbNIhaldLQ",
	})

	redisServer.SetError("REDISDOWN")
	h.UpdateLongUrl(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateUrlUnavailableShortLink(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockUpdateJSONPost(c, handler.UrlUpdateRequest{
		ShortUrl:   "gaia",
		NewLongUrl: "https://youtu.be/UIbNIhaldLQ",
	})

	h.UpdateLongUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectShortUrlSuccess(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "NpHftVNe"
	initialUrl := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.AddParam("shortUrl", shortUrl)

	h.HandleShortUrlRedirect(c)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestRedirectShortUrlRedisFail(t *testing.T) {
	shortUrl := "NpHftVNe"
	initialUrl := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.AddParam("shortUrl", shortUrl)

	redisServer.SetError("REDISDOWN")
	h.HandleShortUrlRedirect(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestRemoveUrlSuccess(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockRemoveJSONPost(c, handler.UrlRemoveRequest{
		ShortUrl: "dyna",
	})

	h.RemoveShortUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveUrlEmptyShortUrl(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockRemoveJSONPost(c, handler.UrlRemoveRequest{
		ShortUrl: "",
	})

	h.RemoveShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveUrlRedisFail(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockRemoveJSONPost(c, handler.UrlRemoveRequest{
		ShortUrl: "dyna",
	})

	redisServer.SetError("REDISDOWN")
	h.RemoveShortUrl(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemoveUrlUnavailableShortLink(t *testing.T) {
	shortener := shortener.NewShortener()
	cfg, err := config.NewConfig("../test.application.yml")
	assert.NoError(t, err)
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	shortUrl := "dyna"
	initialUrl := "https://youtu.be/8LhMu4bQTQU"
	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	storageService := store.StorageService{
		RedisClient: redisClient,
	}
	h := handler.NewHandler(shortener, cfg, &storageService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockRemoveJSONPost(c, handler.UrlRemoveRequest{
		ShortUrl: "gaia",
	})

	h.RemoveShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
