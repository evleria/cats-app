package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitPrice struct {
	channel      *amqp.Channel
	exchangeName string
}

// NewRabbitPriceProducer creates a new producer for rabbitMQ
func NewRabbitPriceProducer(channel *amqp.Channel, exchangeName string) (Price, error) {
	err := channel.ExchangeDeclare(exchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &rabbitPrice{
		channel:      channel,
		exchangeName: exchangeName,
	}, nil
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

	fmt.Printf("producing message to rabbit: {%v, %f}\n", id, price)
	return r.channel.Publish(
		r.exchangeName,
		"",
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
