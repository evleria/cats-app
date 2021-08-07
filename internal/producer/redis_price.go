package producer

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type redisPrice struct {
	redis *redis.Client
}

// NewRedisPriceProducer creates new producer to price stream
func NewRedisPriceProducer(redisClient *redis.Client) Price {
	return &redisPrice{
		redis: redisClient,
	}
}

func (p *redisPrice) Produce(ctx context.Context, id uuid.UUID, price float64) error {
	fmt.Printf("producing message to redis: {%v, %f}\n", id, price)
	args := &redis.XAddArgs{
		Stream: "price",
		Values: map[string]interface{}{
			"id":    id.String(),
			"price": price,
		},
	}
	return p.redis.XAdd(ctx, args).Err()
}
