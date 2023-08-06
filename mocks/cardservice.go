package mocks

import (
	"context"
	"mime/multipart"
	"mtg-report/internal/core/dtos"

	"github.com/stretchr/testify/mock"
)

type CardServiceMock struct {
	mock.Mock
}

func NewCardServiceMock() *CardServiceMock {
	return &CardServiceMock{}
}

func (c *CardServiceMock) InsertCard(ctx context.Context, card dtos.RequestInsertCard) (dtos.ResponseInsertCard, error) {
	args := c.Called(ctx, card)
	return args.Get(0).(dtos.ResponseInsertCard), args.Error(1)
}

func (c *CardServiceMock) InsertCards(ctx context.Context, file multipart.File) (int64, int64) {
	args := c.Called(ctx, file)
	return args.Get(0).(int64), args.Get(1).(int64)
}

func (c *CardServiceMock) GetCardbyID(ctx context.Context, id string) (dtos.ResponseCard, error) {
	args := c.Called(ctx, id)
	return args.Get(0).(dtos.ResponseCard), args.Error(1)
}

func (c *CardServiceMock) GetCards(ctx context.Context, filters map[string]string) ([]dtos.ResponseCard, error) {
	args := c.Called(ctx, filters)
	return args.Get(0).([]dtos.ResponseCard), args.Error(1)
}

func (c *CardServiceMock) DeleteCard(ctx context.Context, id string) error {
	args := c.Called(ctx, id)
	return args.Error(0)
}

func (c *CardServiceMock) GetCardHistory(ctx context.Context, id string) ([]dtos.ResponseCard, error) {
	args := c.Called(ctx, id)
	return args.Get(0).([]dtos.ResponseCard), args.Error(1)
}

func (c *CardServiceMock) UpdateCard(ctx context.Context, cardRequest dtos.RequestUpdateCard) (dtos.ResponseInsertCard, error) {
	args := c.Called(ctx, cardRequest)
	return args.Get(0).(dtos.ResponseInsertCard), args.Error(1)
}
