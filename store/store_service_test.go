package store_test

import (
	"context"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func TestStoreInit(t *testing.T) {
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

func TestInsertionAndRetrieval(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	ctx := context.TODO()

	storageService := store.StorageService{
		RedisClient: redisClient,
	}

	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortUrl := "Jsz4k57oAX"

	storageService.SaveUrlMapping(ctx, shortUrl, initialUrl, userUUId)

	retrievedUrl := storageService.RetrieveInitialUrl(ctx, shortUrl)

	assert.Equal(t, initialUrl, retrievedUrl)
}
