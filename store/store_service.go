package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"source.golabs.io/daniel.santoso/url-blaster/config"
)

type StorageServiceI interface {
	SaveUrlMapping(shortUrl string, originalUrl string, userId string)
	RetrieveInitialUrl(shortUrl string) string
}

type StorageService struct {
	redisClient *redis.Client
	cfg         config.Config
	ctx         context.Context
}

const CacheDuration = 6 * time.Hour

func NewStorageService(cfg config.Config, ctx context.Context) StorageServiceI {
	redisClient := initializeRedis(cfg, ctx)

	return &StorageService{
		redisClient: redisClient,
		cfg:         cfg,
		ctx:         ctx,
	}
}

func initializeRedis(cfg config.Config, ctx context.Context) *redis.Client {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.StoragePort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	return redisClient
}

func (s *StorageService) SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := s.redisClient.Set(s.ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

func (s *StorageService) RetrieveInitialUrl(shortUrl string) string {
	result, err := s.redisClient.Get(s.ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed retrieving inital url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
	return result
}
