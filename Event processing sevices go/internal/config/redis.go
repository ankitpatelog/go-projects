package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password for local
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20, // Redis pool (IMPORTANT)
		MinIdleConns: 10,
	})

	// Redis ping MUST use context
	if err := Redis.Ping(Ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println(" Redis connected")
}
