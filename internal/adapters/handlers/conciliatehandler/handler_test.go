package conciliatehandler

import (
	"context"
	"fmt"
	"testing"

	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	mockPriceService := mocks.NewPriceServiceMock()
	mockLogger := mocks.NewLogMock()

	handler := New(mockPriceService, mockLogger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockPriceService, handler.PriceService)
	assert.Equal(t, mockLogger, handler.log)
}

func TestConciliate_Success(t *testing.T) {
	mockPriceService := mocks.NewPriceServiceMock()
	mockLogger := mocks.NewLogMock()
	mockCustom := mocks.NewCustomMock()

	handler := New(mockPriceService, mockLogger)

	cardsUpdated := int64(5)

	mockLogger.On("Info", mock.Anything).Once()
	mockPriceService.On("Conciliate", mock.Anything).Return(cardsUpdated, nil)
	mockLogger.On("WithFields", mock.AnythingOfType("logrus.Fields")).Return(mockCustom)
	mockCustom.On("Info", mock.Anything).Once()

	err := handler.Conciliate(context.Background())

	assert.NoError(t, err)
	mockPriceService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockCustom.AssertExpectations(t)
}

func TestConciliate_ServiceError(t *testing.T) {
	mockPriceService := mocks.NewPriceServiceMock()
	mockLogger := mocks.NewLogMock()
	mockCustom := mocks.NewCustomMock()

	handler := New(mockPriceService, mockLogger)

	expectedError := fmt.Errorf("service error")
	cardsUpdated := int64(0)

	mockLogger.On("Info", mock.Anything).Once()
	mockPriceService.On("Conciliate", mock.Anything).Return(cardsUpdated, expectedError)
	mockLogger.On("WithError", expectedError).Return(mockCustom)
	mockCustom.On("Error", mock.Anything).Once()
	mockLogger.On("WithFields", mock.AnythingOfType("logrus.Fields")).Return(mockCustom)
	mockCustom.On("Info", mock.Anything).Once()

	err := handler.Conciliate(context.Background())

	// Handler always returns nil, but logs the error
	assert.NoError(t, err)
	mockPriceService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockCustom.AssertExpectations(t)
}
