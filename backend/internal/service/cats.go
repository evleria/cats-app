package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/evleria/mongo-crud/backend/internal/producer"
	"github.com/evleria/mongo-crud/backend/internal/repository"
)

type Cats interface {
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error
	CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error)
}

type cats struct {
	repository    repository.Cats
	priceProducer producer.Price
}

func NewCatsService(catsRepository repository.Cats, priceProducer producer.Price) Cats {
	return &cats{
		repository:    catsRepository,
		priceProducer: priceProducer,
	}
}

func (c *cats) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	err := c.repository.UpdatePrice(ctx, id, price)
	if err != nil {
		return err
	}

	err = c.priceProducer.Produce(ctx, id, price)
	return err
}

func (c *cats) CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error) {
	id, err := c.repository.Insert(ctx, name, color, age, price)
	if err != nil {
		return id, err
	}

	err = c.priceProducer.Produce(ctx, id, price)
	return id, err
}
