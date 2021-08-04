// Package producer provides producing of messages to price stream
package producer

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Price provides producing to price stream
type Price interface {
	Produce(ctx context.Context, id uuid.UUID, price float64) error
}

type price struct {
	redis *redis.Client
}

// NewPriceProducer creates new producer to price stream
func NewPriceProducer(redisClient *redis.Client) Price {
	return &price{
		redis: redisClient,
	}
}

func (p *price) Produce(ctx context.Context, id uuid.UUID, price float64) error {
	args := &redis.XAddArgs{
		Stream: "price",
		Values: map[string]interface{}{
			"id":    id.String(),
			"price": price,
		},
	}
	return p.redis.XAdd(ctx, args).Err()
}
