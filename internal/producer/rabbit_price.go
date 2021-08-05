package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitPrice struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// NewRabbitPriceProducer creates a new producer for rabbitMQ
func NewRabbitPriceProducer(channel *amqp.Channel, exchange, routingKey string) Price {
	return &rabbitPrice{
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

func (r *rabbitPrice) Produce(_ context.Context, id uuid.UUID, price float64) error {
	msg := message{
		ID:    id.String(),
		Price: price,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Printf("producing: %v\n", msg)

	return r.channel.Publish(
		r.exchange,
		r.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
}

type message struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}
