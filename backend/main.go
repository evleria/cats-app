package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/evleria/mongo-crud/backend/internal/handler"
	"github.com/evleria/mongo-crud/backend/internal/repository"
	"github.com/evleria/mongo-crud/backend/internal/service"
	"github.com/evleria/mongo-crud/backend/internal/stream"
)

func main() {
	mongoURI, dbName := getMongoURI()
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	check(err)

	redisClient := redis.NewClient(getRedisOptions())
	_, err = redisClient.Ping(context.Background()).Result()
	check(err)

	catsRepository := repository.NewCatsRepository(mongoClient, dbName)
	priceStream := stream.NewPriceStream(redisClient)
	catsService := service.NewCatsService(catsRepository, priceStream)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	catsGroup := e.Group("/api/cats")
	catsGroup.GET("", handler.GetAllCats(catsRepository))
	catsGroup.GET("/:id", handler.GetCat(catsRepository))
	catsGroup.POST("", handler.AddNewCat(catsService))
	catsGroup.PUT("/:id/price", handler.UpdatePrice(catsService))
	catsGroup.DELETE("/:id", handler.DeleteCat(catsRepository))

	check(e.Start(":5000"))
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
