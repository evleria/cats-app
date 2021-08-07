package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitPrice struct {
	channel   *amqp.Channel
	queueName string
}

// NewRabbitPriceConsumer creates new rabbit price consumer
func NewRabbitPriceConsumer(channel *amqp.Channel, queueName, exchange string) (Price, error) {
	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = channel.QueueBind(q.Name, "", exchange, false, nil)
	if err != nil {
		return nil, err
	}

	return &rabbitPrice{
		channel:   channel,
		queueName: q.Name,
	}, nil
}

func (p *rabbitPrice) Consume(_ context.Context, callbackFunc func(id uuid.UUID, price float64) error) error {
	msgs, err := p.channel.Consume(p.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		id, price, err := decodeRabbitMessage(msg.Body)
		if err != nil {
			return err
		}

		fmt.Printf("consumed message from rabbit: {%v, %f}\n", id, price)
		err = callbackFunc(id, price)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeRabbitMessage(bytes []byte) (id uuid.UUID, price float64, err error) {
	var message struct {
		ID    string  `json:"id"`
		Price float64 `json:"price"`
	}
	err = json.Unmarshal(bytes, &message)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	id, err = uuid.Parse(message.ID)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	return id, message.Price, nil
}
