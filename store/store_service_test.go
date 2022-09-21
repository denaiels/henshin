package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"source.golabs.io/daniel.santoso/url-blaster/config"
	"source.golabs.io/daniel.santoso/url-blaster/store"
)

func TestStoreInit(t *testing.T) {
	cfg, err := config.NewConfig("../dev.application.yml")
	if err != nil {
		return
	}
	ctx := context.Background()
	store := store.NewStorageService(cfg, ctx)
	assert.True(t, store != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	cfg, err := config.NewConfig("../dev.application.yml")
	if err != nil {
		return
	}
	ctx := context.Background()
	store := store.NewStorageService(cfg, ctx)
	initialUrl := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortUrl := "Jsz4k57oAX"

	store.SaveUrlMapping(shortUrl, initialUrl, userUUId)

	retrievedUrl := store.RetrieveInitialUrl(shortUrl)

	assert.Equal(t, initialUrl, retrievedUrl)
}
