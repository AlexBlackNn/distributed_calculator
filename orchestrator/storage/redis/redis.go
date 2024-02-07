package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/config"
	"strconv"
	"time"
)

type Cache struct {
	client *redis.Client
}

func New(cfg *config.Config) *Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Cache{client: redisClient}
}

func (s *Cache) SaveOperation(
	ctx context.Context,
	operation string,
	result float64,
	ttl time.Duration,
) error {
	const op = "DATA LAYER: storage.redis.SaveToken"

	err := s.client.Set(ctx, operation, result, ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Cache) GetOperation(
	ctx context.Context,
	operation string,
) (float64, error) {
	const op = "DATA LAYER: storage.redis.GetToken"

	val, err := s.client.Get(ctx, operation).Result()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	fval, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return fval, nil
}
