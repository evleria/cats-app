package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/evleria/mongo-crud/internal/consumer"
	"github.com/evleria/mongo-crud/internal/handler"
	"github.com/evleria/mongo-crud/internal/producer"
	"github.com/evleria/mongo-crud/internal/repository"
	"github.com/evleria/mongo-crud/internal/service"
)

func main() {
	mongoURI, dbName := getMongoURI()
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	check(err)

	redisClient := redis.NewClient(getRedisOptions())
	_, err = redisClient.Ping(context.Background()).Result()
	check(err)

	catsRepository := repository.NewCatsRepository(mongoClient, dbName)
	priceProducer := producer.NewPriceProducer(redisClient)
	catsService := service.NewCatsService(catsRepository, priceProducer)

	e := echo.New()

	e.Use(middleware.Recover())

	catsGroup := e.Group("/api/cats")
	catsGroup.GET("", handler.GetAllCats(catsRepository))
	catsGroup.GET("/:id", handler.GetCat(catsRepository))
	catsGroup.POST("", handler.AddNewCat(catsService))
	catsGroup.PUT("/:id/price", handler.UpdatePrice(catsService))
	catsGroup.DELETE("/:id", handler.DeleteCat(catsRepository))

	go consumePrices(redisClient)

	check(e.Start(":5000"))
}

func consumePrices(redisClient *redis.Client) {
	priceConsumer := consumer.NewPriceConsumer(redisClient)

	lastID := "0"
	for {
		id, err := priceConsumer.Consume(context.Background(), lastID, func(id uuid.UUID, price float64) error {
			log.Println(id.String(), price)
			return nil
		})

		if err != nil {
			log.Println(err.Error())
		}
		lastID = id
	}
}

func getMongoURI() (mongoURI, dbName string) {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
			getEnvVar("MONGO_USER", "root"),
			getEnvVar("MONGO_PASS", "password"),
			getEnvVar("MONGO_HOST", "localhost"),
			getEnvVar("MONGO_PORT", "27017")),
		getEnvVar("MONGO_DB", "test")
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
