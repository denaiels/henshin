package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"source.golabs.io/daniel.santoso/url-blaster/config"
)

type StorageServiceI interface {
	SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string) error
	RetrieveInitialUrl(ctx context.Context, shortUrl string) (string, error)
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
		log.Panic().Msg(fmt.Sprintf("Error init Redis: %v", err))
	}

	log.Info().Msg(fmt.Sprintf("\nRedis started successfully: pong message = {%s}", pong))
	return redisClient
}

func (s *StorageService) SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string) error {
	err := s.RedisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) RetrieveInitialUrl(ctx context.Context, shortUrl string) (string, error) {
	result, err := s.RedisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
