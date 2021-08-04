package main

import (
	"context"
	"fmt"
	"github.com/evleria/mongo-crud/consumer/internal/consumer"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"os"
)

func main() {
	redisClient := redis.NewClient(getRedisOptions())
	_, err := redisClient.Ping(context.Background()).Result()
	check(err)

	priceConsumer := consumer.NewPriceConsumer(redisClient)

	lastId := "0"
	for {
		id, err := priceConsumer.Consume(context.Background(), lastId, func(id uuid.UUID, price float64) error {
			log.Println(id.String(), price)
			return nil
		})

		if err != nil {
			log.Println(err.Error())
		}
		lastId = id
	}
}

func getRedisOptions() *redis.Options {
	addr := fmt.Sprintf("%s:%s",
		getEnvVar("REDIS_HOST", "localhost"),
		getEnvVar("REDIS_PORT", "6379"),
	)
	pass := getEnvVar("REDIS_PASS", "")

	return &redis.Options{
		Addr:     addr,
		Password: pass,
	}
}

func getEnvVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
