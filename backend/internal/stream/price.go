package stream

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Price interface {
	SendPriceUpdate(ctx context.Context, id uuid.UUID, price float64) error
}

type price struct {
	redis *redis.Client
}

func NewPriceStream(redisClient *redis.Client) Price {
	return &price{
		redis: redisClient,
	}
}

func (p *price) SendPriceUpdate(ctx context.Context, id uuid.UUID, price float64) error {
	args := &redis.XAddArgs{
		Stream: "price",
		Values: map[string]interface{}{
			"id":    id.String(),
			"price": price,
		},
	}
	return p.redis.XAdd(ctx, args).Err()
}
