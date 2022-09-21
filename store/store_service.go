package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"source.golabs.io/daniel.santoso/url-blaster/config"
)

type StorageServiceI interface {
	SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string)
	RetrieveInitialUrl(ctx context.Context, shortUrl string) string
}

type StorageService struct {
	Cfg         *config.Config
	RedisClient *redis.Client
}

const CacheDuration = 6 * time.Hour

func NewStorageService(cfg *config.Config, ctx context.Context) StorageServiceI {
	redisClient := initializeRedis(cfg, ctx)

	return &StorageService{
		Cfg:         cfg,
		RedisClient: redisClient,
	}
}

func initializeRedis(cfg *config.Config, ctx context.Context) *redis.Client {
	address := fmt.Sprintf("%s:%s", cfg.StorageHost, cfg.StoragePort)

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

func (s *StorageService) SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string) {
	err := s.RedisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

func (s *StorageService) RetrieveInitialUrl(ctx context.Context, shortUrl string) string {
	result, err := s.RedisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed retrieving inital url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
	return result
}
