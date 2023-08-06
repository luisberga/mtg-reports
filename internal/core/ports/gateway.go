package ports

import (
	"context"
	"mtg-report/internal/core/domain"
)

type CardGateway interface {
	GetCardPrice(ctx context.Context, card domain.Cards) (float64, error)
}

type ExchangeGateway interface {
	GetUSD(ctx context.Context) (float64, error)
}
