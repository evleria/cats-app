// Package producer provides producing of messages of new price
package producer

import (
	"context"

	"github.com/google/uuid"
)

// Price provides producing to price stream
type Price interface {
	Produce(ctx context.Context, id uuid.UUID, price float64) error
}
