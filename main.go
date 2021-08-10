package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/evleria/cats-app/internal/config"
	"github.com/evleria/cats-app/internal/consumer"
	grpcService "github.com/evleria/cats-app/internal/grpc"
	"github.com/evleria/cats-app/internal/handler"
	"github.com/evleria/cats-app/internal/producer"
	"github.com/evleria/cats-app/internal/repository"
	"github.com/evleria/cats-app/internal/service"
	"github.com/evleria/cats-app/protocol/pb"
)

func main() {
	cfg := new(config.Сonfig)
	check(env.Parse(cfg))

	mongoClient, mongoDB := getMongo(cfg)
	defer mongoClient.Disconnect(context.Background()) //nolint:errcheck,gocritic

	redisClient := getRedis(cfg)
	defer redisClient.Close() //nolint:errcheck,gocritic

	rabbitClient := getRabbit(cfg)
	defer rabbitClient.Close() //nolint:errcheck,gocritic

	rabbitChannel, err := rabbitClient.Channel()
	check(err)
	defer rabbitChannel.Close() //nolint:errcheck,gocritic

	go consumePrices(redisClient, rabbitChannel, cfg.ConsumerNumber)

	catsRepository := repository.NewCatsRepository(mongoDB)
	priceProducer := producer.NewRedisPriceProducer(redisClient)
	catsService := service.NewCatsService(catsRepository, priceProducer)

	e := echo.New()
	e.Use(middleware.Recover())

	catsGroup := e.Group("/api/cats")
	catsGroup.GET("", handler.GetAllCats(catsService))
	catsGroup.GET("/:id", handler.GetCat(catsService))
	catsGroup.POST("", handler.AddNewCat(catsService))
	catsGroup.PUT("/:id/price", handler.UpdatePrice(catsService))
	catsGroup.DELETE("/:id", handler.DeleteCat(catsService))

	go startGrpcServer(catsService, ":6000")

	check(e.Start(":5000"))
}

func startGrpcServer(catsService service.Cats, port string) {
	listener, err := net.Listen("tcp", port)
	check(err)

	s := grpc.NewServer()
	pb.RegisterCatsServiceServer(s, grpcService.NewCatsService(catsService))
	reflection.Register(s)

	fmt.Printf("Starting gRPC server on port %s\n", port)
	check(s.Serve(listener))
}

func consumePrices(redisClient *redis.Client, rabbitChannel *amqp.Channel, consumerNumber int) {
	queueName := fmt.Sprintf("price_%d", consumerNumber)
	rabbitPriceProducer, err := producer.NewRabbitPriceProducer(rabbitChannel, "price")
	check(err)
	rabbitPriceConsumer, err := consumer.NewRabbitPriceConsumer(rabbitChannel, queueName, "price")
	check(err)

	redisPriceConsumer := consumer.NewRedisPriceConsumer(redisClient, fmt.Sprintf("%d000-0", time.Now().Unix()))
	go func() {
		err := rabbitPriceConsumer.Consume(context.Background(), func(id uuid.UUID, price float64) error {
			return nil
		})
		check(err)
	}()

	go func() {
		err := redisPriceConsumer.Consume(context.Background(), func(id uuid.UUID, price float64) error {
			err := rabbitPriceProducer.Produce(context.Background(), id, price)
			if err != nil {
				log.Println(err.Error())
			}
			return err
		})
		check(err)
	}()
}

func getMongo(cfg *config.Сonfig) (*mongo.Client, *mongo.Database) {
	mongoURI, dbName := getMongoURI(cfg)

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	check(err)

	db := mongoClient.Database(dbName)
	return mongoClient, db
}

func getMongoURI(cfg *config.Сonfig) (mongoURI, dbName string) {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d",
			cfg.MongoUser,
			cfg.MongoPassword,
			cfg.MongoHost,
			cfg.MongoPort),
		cfg.MongoDB
}

func getRedis(cfg *config.Сonfig) *redis.Client {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
	}

	redisClient := redis.NewClient(opts)
	_, err := redisClient.Ping(context.Background()).Result()
	check(err)

	return redisClient
}

func getRabbit(cfg *config.Сonfig) *amqp.Connection {
	connection, err := amqp.Dial(getRabbitURL(cfg))
	check(err)
	return connection
}

func getRabbitURL(cfg *config.Сonfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitUser,
		cfg.RabbitPass,
		cfg.RabbitHost,
		cfg.RabbitPort,
	)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
