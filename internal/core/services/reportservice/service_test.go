package reportservice

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

func TestNew(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.ReportRepository)
	assert.Equal(t, mockEmail, service.Email)
	assert.Equal(t, mockLogger, service.log)
}

func TestProcessAndSend_Success(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	now := time.Now()
	expectedCards := []domain.Cards{
		{
			ID:              1,
			Name:            "Lightning Bolt",
			SetName:         "Alpha",
			CollectorNumber: "161",
			Foil:            false,
			CardsDetails: domain.CardsDetails{
				LastPrice:   10.50,
				OldPrice:    9.00,
				PriceChange: 1.50,
				LastUpdate:  &now,
			},
		},
	}

	expectedPrice := domain.CardsPrice{
		OldPrice:    100.00,
		NewPrice:    110.50,
		PriceChange: 10.50,
		LastUpdate:  &now,
	}

	mockRepo.On("InsertTotalPrice", mock.Anything).Return(nil)
	mockRepo.On("GetCardsReport", mock.Anything).Return(expectedCards, nil)
	mockRepo.On("GetTotalPrice", mock.Anything).Return(expectedPrice, nil)
	mockEmail.On("SendEmail", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := service.ProcessAndSend(context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestProcessAndSend_InsertTotalPriceError(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	mockRepo.On("InsertTotalPrice", mock.Anything).Return(fmt.Errorf("database error"))

	err := service.ProcessAndSend(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service failed to insert total price in process and send")
	mockRepo.AssertExpectations(t)
}

func TestProcessAndSend_GetCardsReportError(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	mockRepo.On("InsertTotalPrice", mock.Anything).Return(nil)
	mockRepo.On("GetCardsReport", mock.Anything).Return([]domain.Cards(nil), fmt.Errorf("query error"))

	err := service.ProcessAndSend(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service failed to get cards reports in process and send")
	mockRepo.AssertExpectations(t)
}

func TestProcessAndSend_GetTotalPriceError(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	expectedCards := []domain.Cards{}

	mockRepo.On("InsertTotalPrice", mock.Anything).Return(nil)
	mockRepo.On("GetCardsReport", mock.Anything).Return(expectedCards, nil)
	mockRepo.On("GetTotalPrice", mock.Anything).Return(domain.CardsPrice{}, fmt.Errorf("query error"))

	err := service.ProcessAndSend(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service failed to get total price in process and send")
	mockRepo.AssertExpectations(t)
}

func TestProcessAndSend_SendEmailError(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	expectedCards := []domain.Cards{}
	expectedPrice := domain.CardsPrice{}

	mockRepo.On("InsertTotalPrice", mock.Anything).Return(nil)
	mockRepo.On("GetCardsReport", mock.Anything).Return(expectedCards, nil)
	mockRepo.On("GetTotalPrice", mock.Anything).Return(expectedPrice, nil)
	mockEmail.On("SendEmail", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(fmt.Errorf("email error"))

	err := service.ProcessAndSend(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service failed to send email in process and send")
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestFormatCardsTable(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	now := time.Now()
	cards := []domain.Cards{
		{
			ID:              1,
			Name:            "Lightning Bolt",
			SetName:         "Alpha",
			CollectorNumber: "161",
			Foil:            false,
			CardsDetails: domain.CardsDetails{
				LastPrice:   10.50,
				OldPrice:    9.00,
				PriceChange: 1.50, // Positive change - green
				LastUpdate:  &now,
			},
		},
		{
			ID:              2,
			Name:            "Counterspell",
			SetName:         "Alpha",
			CollectorNumber: "50",
			Foil:            true,
			CardsDetails: domain.CardsDetails{
				LastPrice:   4.00,
				OldPrice:    5.25,
				PriceChange: -1.25, // Negative change - red
				LastUpdate:  &now,
			},
		},
		{
			ID:              3,
			Name:            "Black Lotus",
			SetName:         "Alpha",
			CollectorNumber: "232",
			Foil:            false,
			CardsDetails: domain.CardsDetails{
				LastPrice:   5000.00,
				OldPrice:    5000.00,
				PriceChange: 0.00, // No change - black
				LastUpdate:  nil,  // Test nil timestamp
			},
		},
	}

	result := service.formatCardsTable(cards)

	// Check that result contains expected elements
	assert.Contains(t, result, "<tr>")
	assert.Contains(t, result, "<td")
	assert.Contains(t, result, "Lightning Bolt")
	assert.Contains(t, result, "Counterspell")
	assert.Contains(t, result, "Black Lotus")
	assert.Contains(t, result, "Alpha")
	assert.Contains(t, result, "161")
	assert.Contains(t, result, "50")
	assert.Contains(t, result, "232")
	assert.Contains(t, result, "10.50")
	assert.Contains(t, result, "4.00")
	assert.Contains(t, result, "5000.00")
	assert.Contains(t, result, "color: green") // Positive change
	assert.Contains(t, result, "color: red")   // Negative change
	assert.Contains(t, result, "color: black") // No change
	assert.Contains(t, result, "</table>")
}

func TestFormatCardsPrice(t *testing.T) {
	mockRepo := mocks.NewReportRepositoryMock()
	mockEmail := mocks.NewEmailMock()
	mockLogger := mocks.NewLogMock()

	service := New(mockRepo, mockEmail, mockLogger)

	tests := []struct {
		name     string
		price    domain.CardsPrice
		expected string
	}{
		{
			name: "Positive change",
			price: domain.CardsPrice{
				OldPrice:    100.00,
				NewPrice:    110.50,
				PriceChange: 10.50,
			},
			expected: "increased",
		},
		{
			name: "Negative change",
			price: domain.CardsPrice{
				OldPrice:    100.00,
				NewPrice:    85.25,
				PriceChange: -14.75,
			},
			expected: "decreased",
		},
		{
			name: "No change",
			price: domain.CardsPrice{
				OldPrice:    100.00,
				NewPrice:    100.00,
				PriceChange: 0.00,
			},
			expected: "stayed the same",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.formatCardsPrice(tt.price)

			assert.Contains(t, result, tt.expected)
			assert.Contains(t, result, fmt.Sprintf("R$%.2f", tt.price.OldPrice))
			assert.Contains(t, result, fmt.Sprintf("R$%.2f", tt.price.NewPrice))
			assert.Contains(t, result, fmt.Sprintf("R$%.2f", tt.price.PriceChange))
		})
	}
}
