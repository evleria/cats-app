package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/streadway/amqp"
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
	defer mongoClient.Disconnect(context.Background()) //nolint:errcheck,gocritic

	redisClient := redis.NewClient(getRedisOptions())
	_, err = redisClient.Ping(context.Background()).Result()
	check(err)
	defer redisClient.Close() //nolint:errcheck,gocritic

	rabbitClient, err := amqp.Dial(getRabbitURL())
	check(err)
	defer rabbitClient.Close() //nolint:errcheck,gocritic

	rabbitChannel, err := rabbitClient.Channel()
	check(err)
	defer rabbitChannel.Close() //nolint:errcheck,gocritic

	_, err = rabbitChannel.QueueDeclare(
		"price",
		true,
		false,
		false,
		false,
		nil)
	check(err)

	go consumePrices(redisClient, rabbitChannel)

	catsRepository := repository.NewCatsRepository(mongoClient, dbName)
	priceProducer := producer.NewRedisPriceProducer(redisClient)
	catsService := service.NewCatsService(catsRepository, priceProducer)

	e := echo.New()
	e.Use(middleware.Recover())

	catsGroup := e.Group("/api/cats")
	catsGroup.GET("", handler.GetAllCats(catsRepository))
	catsGroup.GET("/:id", handler.GetCat(catsRepository))
	catsGroup.POST("", handler.AddNewCat(catsService))
	catsGroup.PUT("/:id/price", handler.UpdatePrice(catsService))
	catsGroup.DELETE("/:id", handler.DeleteCat(catsRepository))

	check(e.Start(":5000"))
}

func consumePrices(redisClient *redis.Client, rabbitChannel *amqp.Channel) {
	redisPriceConsumer := consumer.NewPriceConsumer(redisClient, fmt.Sprintf("%d000-0", time.Now().Unix()))
	rabbitPriceProducer := producer.NewRabbitPriceProducer(rabbitChannel, "", "price")

	for {
		err := redisPriceConsumer.Consume(context.Background(), func(id uuid.UUID, price float64) error {
			return rabbitPriceProducer.Produce(context.Background(), id, price)
		})

		if err != nil {
			log.Println(err.Error())
		}
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

func getRabbitURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		getEnvVar("RABBIT_USER", "guest"),
		getEnvVar("RABBIT_PASS", "guest"),
		getEnvVar("RABBIT_HOST", "localhost"),
		getEnvVar("RABBIT_PORT", "5672"),
	)
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
