package store_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

const CacheDuration = 6 * time.Hour

func TestStoreInitSuccess(t *testing.T) {
	cfg, err := config.NewConfig("../test.application.yml")
	if err != nil {
		return
	}

	redisServer := miniredis.RunT(t)
	address := strings.Split(redisServer.Addr(), ":")
	cfg.StorageHost = address[0]
	cfg.StoragePort = address[1]
	ctx := context.Background()

	store := store.NewStorageService(cfg, ctx)

	assert.True(t, store != nil)
}

func TestSaveUrlMappingSuccess(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	err := storageService.SaveUrlMapping(ctx, shortUrl, initialUrl)
	assert.NoError(t, err)
}

func TestSaveUrlMappingFail(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisServer.SetError("REDISDOWN")
	err := storageService.SaveUrlMapping(ctx, shortUrl, initialUrl)
	assert.Error(t, err)
}

func TestRetrieveInitialUrlSuccess(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	retrievedUrl, err := storageService.RetrieveInitialUrl(ctx, shortUrl)
	assert.Equal(t, initialUrl, retrievedUrl)
	assert.NoError(t, err)
}

func TestRetrieveInitialUrlFail(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	redisServer.SetError("REDISDOWN")
	retrievedUrl, err := storageService.RetrieveInitialUrl(ctx, shortUrl)
	assert.NotEqual(t, initialUrl, retrievedUrl)
	assert.Error(t, err)
}

func TestCheckIfShortUrlExistsSuccess(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	shortUrlExists := storageService.CheckIfShortUrlExists(ctx, "Jsz4k57oAX")

	assert.True(t, shortUrlExists)
}

func TestCheckIfShortUrlExistsFail(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	shortUrlExists := storageService.CheckIfShortUrlExists(ctx, "awokawok")

	assert.False(t, shortUrlExists)
}

func TestDeleteUrlMappingSuccess(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	err := storageService.DeleteUrlMapping(ctx, "Jsz4k57oAX")
	assert.NoError(t, err)
}

func TestDeleteUrlMappingRedisFail(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortUrl := "Jsz4k57oAX"

	redisClient.Set(ctx, shortUrl, initialUrl, CacheDuration)

	redisServer.SetError("REDISDOWN")
	err := storageService.DeleteUrlMapping(ctx, "Jsz4k57oAX")
	assert.Error(t, err)
}
