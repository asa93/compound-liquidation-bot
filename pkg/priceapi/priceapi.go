package priceapi

import (
	"context"
)

// PriceAPI represents crypto price API
type PriceAPI interface {
	// Get name returns crypto API name
	GetName() string
	// GetPrice returns crypto price
	GetPrice(ctx context.Context) (float64, error)
}
