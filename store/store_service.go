package store

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"source.golabs.io/daniel.santoso/url-blaster/config"
)

type StorageServiceI interface {
	SaveUrlMapping(ctx context.Context, shortUrl, originalUrl string) error
	CheckIfShortUrlExists(ctx context.Context, shortUrl string) bool
	RetrieveInitialUrl(ctx context.Context, shortUrl string) (string, error)
	DeleteUrlMapping(ctx context.Context, shortUrl string) error
}

type StorageService struct {
	Cfg         *config.Config
	RedisClient *redis.Client
}

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
		log.Fatal().Msg(fmt.Sprintf("Error init Redis: %v", err))
	}

	log.Info().Msg(fmt.Sprintf("\nRedis started successfully: pong message = {%s}", pong))
	return redisClient
}

func (s *StorageService) SaveUrlMapping(ctx context.Context, shortUrl, originalUrl string) error {
	err := s.RedisClient.Set(ctx, shortUrl, originalUrl, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) CheckIfShortUrlExists(ctx context.Context, shortUrl string) bool {
	_, err := s.RedisClient.Get(ctx, shortUrl).Result()
	return err != redis.Nil
}

func (s *StorageService) RetrieveInitialUrl(ctx context.Context, shortUrl string) (string, error) {
	result, err := s.RedisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *StorageService) DeleteUrlMapping(ctx context.Context, shortUrl string) error {
	err := s.RedisClient.Del(ctx, shortUrl).Err()
	if err != nil {
		return err
	}

	return nil
}
