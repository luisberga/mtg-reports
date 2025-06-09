package conciliaterepo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"mtg-report/internal/core/domain"
	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertCardDetails_Success(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockResult := mocks.NewResultMock()

	repo := New(mockDB)

	now := time.Now()
	cardDetails := []domain.CardsDetails{
		{
			CardID:      1,
			LastPrice:   10.50,
			OldPrice:    9.00,
			PriceChange: 1.50,
			LastUpdate:  &now,
		},
		{
			CardID:      2,
			LastPrice:   5.25,
			OldPrice:    4.00,
			PriceChange: 1.25,
			LastUpdate:  &now,
		},
	}

	mockResult.On("RowsAffected").Return(int64(2), nil)
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, nil)

	err := repo.InsertCardDetails(context.Background(), cardDetails)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestInsertCardDetails_EmptySlice(t *testing.T) {
	mockDB := mocks.NewClientMock()

	repo := New(mockDB)

	cardDetails := []domain.CardsDetails{}

	err := repo.InsertCardDetails(context.Background(), cardDetails)

	assert.NoError(t, err)
	// Não deve haver chamadas ao banco para slice vazio
	mockDB.AssertNotCalled(t, "ExecContext")
}

func TestInsertCardDetails_DatabaseError(t *testing.T) {
	mockDB := mocks.NewClientMock()

	repo := New(mockDB)

	now := time.Now()
	cardDetails := []domain.CardsDetails{
		{
			CardID:      1,
			LastPrice:   10.50,
			OldPrice:    9.00,
			PriceChange: 1.50,
			LastUpdate:  &now,
		},
	}

	mockResult := mocks.NewResultMock()
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, fmt.Errorf("database error"))

	err := repo.InsertCardDetails(context.Background(), cardDetails)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to execute insert statement")
	mockDB.AssertExpectations(t)
}

func TestInsertCardDetails_NoRowsAffected(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockResult := mocks.NewResultMock()

	repo := New(mockDB)

	now := time.Now()
	cardDetails := []domain.CardsDetails{
		{
			CardID:      1,
			LastPrice:   10.50,
			OldPrice:    9.00,
			PriceChange: 1.50,
			LastUpdate:  &now,
		},
	}

	mockResult.On("RowsAffected").Return(int64(0), nil)
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, nil)

	err := repo.InsertCardDetails(context.Background(), cardDetails)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository insert card details failed")
	mockDB.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestGetCardsForUpdate_Success(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	// Simulating that it finds 1 card
	mockRowsScanner.On("Next").Return(true).Once()
	mockRowsScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRowsScanner.On("Next").Return(false).Once()
	mockRowsScanner.On("Err").Return(nil)
	mockRowsScanner.On("Close").Return(nil)

	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsForUpdate(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}

func TestGetCardsForUpdate_NoCards(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	mockRowsScanner.On("Next").Return(false)
	mockRowsScanner.On("Err").Return(nil)
	mockRowsScanner.On("Close").Return(nil)

	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsForUpdate(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Len(t, cards, 0)
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}

func TestGetCardsForUpdate_DatabaseError(t *testing.T) {
	mockDB := mocks.NewClientMock()

	repo := New(mockDB)

	mockRowsScanner := mocks.NewRowsScannerMock()
	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, fmt.Errorf("database error"))

	cards, err := repo.GetCardsForUpdate(context.Background(), 0, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to query in get cards for update")
	assert.Nil(t, cards)
	mockDB.AssertExpectations(t)
}

func TestGetCardsForUpdate_ScanError(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	mockRowsScanner.On("Next").Return(true).Once()
	mockRowsScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("scan error"))
	mockRowsScanner.On("Close").Return(nil)

	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsForUpdate(context.Background(), 0, 10)

	// O erro será causado pela conversão de tipo do ID
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to scan rows in get cards for update")
	assert.Nil(t, cards)
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}

func TestGetCardsForUpdate_WithOffsetAndLimit(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	mockRowsScanner.On("Next").Return(true).Once()
	mockRowsScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRowsScanner.On("Next").Return(false).Once()
	mockRowsScanner.On("Err").Return(nil)
	mockRowsScanner.On("Close").Return(nil)

	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsForUpdate(context.Background(), 10, 5)

	assert.NoError(t, err)
	assert.Len(t, cards, 1) // With zero values due to simple mock
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}
