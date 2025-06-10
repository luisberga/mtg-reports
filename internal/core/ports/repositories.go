package ports

import (
	"context"
	"mtg-report/internal/core/domain"
)

type CardsRepository interface {
	InsertCard(ctx context.Context, card domain.Cards) (domain.Cards, error)
	InsertCards(ctx context.Context, cards []domain.Cards) error
	GetCardbyID(ctx context.Context, id string) (domain.Cards, error)
	GetCards(ctx context.Context, filters map[string]string) ([]domain.Cards, error)
	GetCardsPaginated(ctx context.Context, filters map[string]string, offset, limit int) ([]domain.Cards, error)
	GetCardsCount(ctx context.Context, filters map[string]string) (int64, error)
	DeleteCard(ctx context.Context, id string) error
	GetCardHistory(ctx context.Context, id string) ([]domain.Cards, error)
	GetCardHistoryPaginated(ctx context.Context, id string, offset, limit int) ([]domain.Cards, error)
	GetCardHistoryCount(ctx context.Context, id string) (int64, error)
	UpdateCard(ctx context.Context, card domain.UpdateCard) (domain.Cards, error)
	GetCollectionStats(ctx context.Context) (domain.CollectionStats, error)
}

type ConciliateRepository interface {
	GetCardsForUpdate(ctx context.Context, offset int, limit int) ([]domain.Cards, error)
	InsertCardDetails(ctx context.Context, cards []domain.CardsDetails) error
}

type ReportRepository interface {
	InsertTotalPrice(ctx context.Context) error
	GetCardsReport(ctx context.Context) ([]domain.Cards, error)
	GetTotalPrice(ctx context.Context) (domain.CardsPrice, error)
}
