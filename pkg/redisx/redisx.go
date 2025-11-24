package redisx

import (
	"context"
	"fmt"
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

// GetAddr формирует адрес Redis из конфигурации
func GetAddr() string {
	cfg := config.G()
	return fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
}

// Open открывает соединение с Redis
func Open() {
	cfg := config.G()

	Client = redis.NewClient(&redis.Options{
		Addr:     GetAddr(),
		Password: cfg.Redis.Password,
		DB:       0, // default database
	})

	// Проверяем соединение
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Printf("connected to redis at %s", GetAddr())
}

// Close закрывает соединение с Redis
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// G возвращает глобальный клиент Redis (singleton pattern)
func G() *redis.Client {
	if Client == nil {
		log.Panic("redis client is nil")
	}
	return Client
}
