package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/evleria/mongo-crud/backend/internal/repository"
	"github.com/evleria/mongo-crud/backend/internal/stream"
)

type Cats interface {
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error
	CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error)
}

type cats struct {
	repository  repository.Cats
	priceStream stream.Price
}

func NewCatsService(catsRepository repository.Cats, priceStream stream.Price) Cats {
	return &cats{
		repository:  catsRepository,
		priceStream: priceStream,
	}
}

func (c *cats) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	err := c.repository.UpdatePrice(ctx, id, price)
	if err != nil {
		return err
	}

	return c.priceStream.SendPriceUpdate(ctx, id, price)
}

func (c *cats) CreateNew(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error) {
	id, err := c.repository.Insert(ctx, name, color, age, price)
	if err != nil {
		return id, err
	}

	return id, c.priceStream.SendPriceUpdate(ctx, id, price)
}
