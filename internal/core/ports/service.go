package ports

import (
	"context"
	"mime/multipart"
	"mtg-report/internal/core/dtos"
)

type CardService interface {
	InsertCard(ctx context.Context, card dtos.RequestInsertCard) (dtos.ResponseInsertCard, error)
	InsertCards(ctx context.Context, file multipart.File) (int64, int64)
	GetCardbyID(ctx context.Context, id string) (dtos.ResponseCard, error)
	GetCards(ctx context.Context, filters map[string]string) ([]dtos.ResponseCard, error)
	GetCardsPaginated(ctx context.Context, filters map[string]string, page, limit int) (dtos.ResponsePaginatedCards, error)
	DeleteCard(ctx context.Context, id string) error
	GetCardHistory(ctx context.Context, id string) ([]dtos.ResponseCard, error)
	UpdateCard(ctx context.Context, cardRequest dtos.RequestUpdateCard) (dtos.ResponseInsertCard, error)
}

type PriceService interface {
	Conciliate(ctx context.Context) (int64, error)
}

type ReportService interface {
	ProcessAndSend(ctx context.Context) error
}
