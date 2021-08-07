// Package consumer provides consuming of messages of price
package consumer

import (
	"context"

	"github.com/google/uuid"
)

// Price consuming price messages
type Price interface {
	Consume(ctx context.Context, callbackFunc func(id uuid.UUID, price float64) error) error
}
