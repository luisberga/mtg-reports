package conciliateservice

import (
	"context"
	"fmt"
	"testing"

	"mtg-report/internal/core/domain"
	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	mockConciliateRepo := mocks.NewConciliateRepositoryMock()
	mockCardGateway := mocks.NewCardGatewayMock()
	mockExchangeGateway := mocks.NewExchangeGatewayMock()
	mockLogger := mocks.NewLogMock()
	commitSize := 10

	service := New(mockConciliateRepo, mockCardGateway, mockExchangeGateway, commitSize, mockLogger)

	assert.NotNil(t, service)
	assert.Equal(t, mockConciliateRepo, service.ConciliateRepository)
	assert.Equal(t, mockCardGateway, service.cardGateway)
	assert.Equal(t, mockExchangeGateway, service.exchangegateway)
	assert.Equal(t, commitSize, service.commitSize)
	assert.Equal(t, mockLogger, service.log)
}

func TestConciliate_NoCardsToUpdate(t *testing.T) {
	mockConciliateRepo := mocks.NewConciliateRepositoryMock()
	mockCardGateway := mocks.NewCardGatewayMock()
	mockExchangeGateway := mocks.NewExchangeGatewayMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockConciliateRepo, mockCardGateway, mockExchangeGateway, 10, mockLogger)

	// Mock exchange rate
	mockExchangeGateway.On("GetUSD", mock.Anything).Return(5.0, nil)

	// Mock no cards to update
	mockConciliateRepo.On("GetCardsForUpdate", mock.Anything, 0, 10).Return([]domain.Cards{}, nil)

	// Mock logger calls
	mockLogger.On("Info", mock.Anything).Maybe()

	cardsUpdated, err := service.Conciliate(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(0), cardsUpdated)
	mockExchangeGateway.AssertExpectations(t)
	mockConciliateRepo.AssertExpectations(t)
}

func TestConciliate_ExchangeGatewayError(t *testing.T) {
	mockConciliateRepo := mocks.NewConciliateRepositoryMock()
	mockCardGateway := mocks.NewCardGatewayMock()
	mockExchangeGateway := mocks.NewExchangeGatewayMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockConciliateRepo, mockCardGateway, mockExchangeGateway, 10, mockLogger)

	// Mock exchange rate error - should use default value
	mockExchangeGateway.On("GetUSD", mock.Anything).Return(0.0, fmt.Errorf("exchange error"))

	// Mock no cards to update
	mockConciliateRepo.On("GetCardsForUpdate", mock.Anything, 0, 10).Return([]domain.Cards{}, nil)

	// Mock logger calls for error
	mockLogger.On("Error", mock.Anything).Once()
	mockLogger.On("Info", mock.Anything).Maybe()

	cardsUpdated, err := service.Conciliate(context.Background())

	// Should still work with default exchange rate
	assert.NoError(t, err)
	assert.Equal(t, int64(0), cardsUpdated)
	mockExchangeGateway.AssertExpectations(t)
	mockConciliateRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestLogError(t *testing.T) {
	mockConciliateRepo := mocks.NewConciliateRepositoryMock()
	mockCardGateway := mocks.NewCardGatewayMock()
	mockExchangeGateway := mocks.NewExchangeGatewayMock()
	mockLogger := mocks.NewLogMock()
	mockCustom := mocks.NewCustomMock()

	service := New(mockConciliateRepo, mockCardGateway, mockExchangeGateway, 10, mockLogger)

	card := domain.Cards{
		ID:              1,
		Name:            "Lightning Bolt",
		SetName:         "Alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	testError := fmt.Errorf("test error")

	// Mock logger calls
	mockLogger.On("WithFields", mock.AnythingOfType("logrus.Fields")).Return(mockCustom)
	mockCustom.On("Warn", mock.Anything).Once()

	// Call the private method through reflection or make it public for testing
	// For now, we'll test it indirectly through the public interface
	service.logError(card, testError)

	mockLogger.AssertExpectations(t)
	mockCustom.AssertExpectations(t)
}
