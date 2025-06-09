package reportrepo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"mtg-report/internal/core/domain"
	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertTotalPrice_Success(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockResult := mocks.NewResultMock()

	repo := New(mockDB)

	mockResult.On("RowsAffected").Return(int64(1), nil)
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, nil)

	err := repo.InsertTotalPrice(context.Background())

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestInsertTotalPrice_DatabaseError(t *testing.T) {
	mockDB := mocks.NewClientMock()

	repo := New(mockDB)

	mockResult := mocks.NewResultMock()
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, fmt.Errorf("database error"))

	err := repo.InsertTotalPrice(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to exec insert query in insert total price")
	mockDB.AssertExpectations(t)
}

func TestInsertTotalPrice_NoRowsAffected(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockResult := mocks.NewResultMock()

	repo := New(mockDB)

	mockResult.On("RowsAffected").Return(int64(0), nil)
	mockDB.On("ExecContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockResult, nil)

	err := repo.InsertTotalPrice(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository insert total price failed")
	mockDB.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestGetCardsReport_Success(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	// Simulando que encontra cards
	mockRowsScanner.On("Next").Return(true).Once()
	mockRowsScanner.On("Scan", mock.Anything).Return(nil).Once()
	mockRowsScanner.On("Next").Return(false).Once()

	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsReport(context.Background())

	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}

func TestGetCardsReport_NoCards(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowsScanner := mocks.NewRowsScannerMock()

	repo := New(mockDB)

	mockRowsScanner.On("Next").Return(false)
	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, nil)

	cards, err := repo.GetCardsReport(context.Background())

	assert.Error(t, err)
	assert.IsType(t, domain.ErrCardNotFound{}, err)
	assert.Nil(t, cards)
	mockDB.AssertExpectations(t)
	mockRowsScanner.AssertExpectations(t)
}

func TestGetCardsReport_DatabaseError(t *testing.T) {
	mockDB := mocks.NewClientMock()

	repo := New(mockDB)

	mockRowsScanner := mocks.NewRowsScannerMock()
	mockDB.On("QueryContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowsScanner, fmt.Errorf("database error"))

	cards, err := repo.GetCardsReport(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to exec query in get cards report")
	assert.Nil(t, cards)
	mockDB.AssertExpectations(t)
}

func TestGetTotalPrice_Success(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowScanner := mocks.NewRowScannerMock()

	repo := New(mockDB)

	mockRowScanner.On("Scan").Return(nil)
	mockDB.On("QueryRowContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowScanner)

	result, err := repo.GetTotalPrice(context.Background())

	// With zero values due to mock limitation
	assert.NoError(t, err)
	assert.Equal(t, 0.0, result.OldPrice)
	assert.Equal(t, 0.0, result.NewPrice)
	assert.Equal(t, 0.0, result.PriceChange)
	mockDB.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

func TestGetTotalPrice_NotFound(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowScanner := mocks.NewRowScannerMock()

	repo := New(mockDB)

	mockRowScanner.On("Scan").Return(sql.ErrNoRows)
	mockDB.On("QueryRowContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowScanner)

	_, err := repo.GetTotalPrice(context.Background())

	assert.Error(t, err)
	assert.IsType(t, domain.ErrCardNotFound{}, err)
	mockDB.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

func TestGetTotalPrice_ScanError(t *testing.T) {
	mockDB := mocks.NewClientMock()
	mockRowScanner := mocks.NewRowScannerMock()

	repo := New(mockDB)

	mockRowScanner.On("Scan").Return(fmt.Errorf("scan error"))
	mockDB.On("QueryRowContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowScanner)

	_, err := repo.GetTotalPrice(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository failed to scan row in get total price")
	mockDB.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}
