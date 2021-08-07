package consumer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type redisPrice struct {
	redis  *redis.Client
	lastID string
}

// NewRedisPriceConsumer creates new redis price consumer
func NewRedisPriceConsumer(redisClient *redis.Client, startID string) Price {
	return &redisPrice{
		redis:  redisClient,
		lastID: startID,
	}
}

func (p *redisPrice) Consume(ctx context.Context, callbackFunc func(id uuid.UUID, price float64) error) error {
	for {
		args := &redis.XReadArgs{
			Streams: []string{"price", p.lastID},
		}
		r, err := p.redis.XRead(ctx, args).Result()
		if err != nil {
			return err
		}

		for _, message := range r[0].Messages {
			id, price, err := decodeRedisMessage(message)
			if err != nil {
				return err
			}

			fmt.Printf("consumed message from redis: {%v, %f}\n", id, price)
			err = callbackFunc(id, price)
			if err != nil {
				return err
			}

			p.lastID = message.ID
		}
	}
}

func decodeRedisMessage(message redis.XMessage) (id uuid.UUID, price float64, err error) {
	idStr, ok := message.Values["id"].(string)
	if !ok {
		return id, price, errors.New("cannot convert id to string")
	}
	priceStr, ok := message.Values["price"].(string)
	if !ok {
		return id, price, errors.New("cannot convert price to string")
	}

	id, err = uuid.Parse(idStr)
	if err != nil {
		return id, price, err
	}
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return id, price, err
	}

	return id, price, nil
}
