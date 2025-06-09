package mocks

import (
	"context"
	"mtg-report/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type ReportRepositoryMock struct {
	mock.Mock
}

func NewReportRepositoryMock() *ReportRepositoryMock {
	return &ReportRepositoryMock{}
}

func (m *ReportRepositoryMock) InsertTotalPrice(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *ReportRepositoryMock) GetCardsReport(ctx context.Context) ([]domain.Cards, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Cards), args.Error(1)
}

func (m *ReportRepositoryMock) GetTotalPrice(ctx context.Context) (domain.CardsPrice, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.CardsPrice), args.Error(1)
}
