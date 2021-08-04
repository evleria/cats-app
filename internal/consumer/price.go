// Package consumer provides consuming of messages from price stream
package consumer

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Price consuming for price stream
type Price interface {
	Consume(ctx context.Context, lastID string, callbackFunc func(id uuid.UUID, price float64) error) (string, error)
}

type price struct {
	redis *redis.Client
}

// NewPriceConsumer creates new price consumer
func NewPriceConsumer(redisClient *redis.Client) Price {
	return &price{
		redis: redisClient,
	}
}

func (p *price) Consume(ctx context.Context, lastID string, callbackFunc func(id uuid.UUID, price float64) error) (string, error) {
	args := &redis.XReadArgs{
		Streams: []string{"price", lastID},
	}
	r, err := p.redis.XRead(ctx, args).Result()
	if err != nil {
		return lastID, err
	}

	for _, message := range r[0].Messages {
		id, price, err := decodeMessage(message)
		if err != nil {
			return lastID, err
		}

		err = callbackFunc(id, price)
		if err != nil {
			return lastID, err
		}

		lastID = message.ID
	}

	return lastID, nil
}

func decodeMessage(message redis.XMessage) (id uuid.UUID, price float64, err error) {
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