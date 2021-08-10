// Package service encapsulates usecases
package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/evleria/mongo-crud/internal/producer"
	"github.com/evleria/mongo-crud/internal/repository"
	"github.com/evleria/mongo-crud/internal/repository/entities"
)

// Cats contains usecase logic for cats
type Cats interface {
	GetAll(ctx context.Context) ([]entities.Cat, error)
	GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error)
	CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error
}

type cats struct {
	repository    repository.Cats
	priceProducer producer.Price
}

// NewCatsService creates new cats service
func NewCatsService(catsRepository repository.Cats, priceProducer producer.Price) Cats {
	return &cats{
		repository:    catsRepository,
		priceProducer: priceProducer,
	}
}

func (c *cats) GetAll(ctx context.Context) ([]entities.Cat, error) {
	return c.repository.GetAll(ctx)
}

func (c *cats) GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error) {
	return c.repository.GetOne(ctx, id)
}

func (c *cats) CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error) {
	id, err := c.repository.Insert(ctx, name, color, age, price)
	if err != nil {
		return id, err
	}

	err = c.priceProducer.Produce(ctx, id, price)
	return id, err
}

func (c *cats) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repository.Delete(ctx, id)
}

func (c *cats) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	err := c.repository.UpdatePrice(ctx, id, price)
	if err != nil {
		return err
	}

	err = c.priceProducer.Produce(ctx, id, price)
	return err
}
