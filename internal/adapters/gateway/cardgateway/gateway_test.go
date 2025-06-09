package cardgateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"mtg-report/internal/core/domain"
	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()

	gateway := New(mockWeb, mockLogger)

	assert.NotNil(t, gateway)
	assert.Equal(t, mockWeb, gateway.web)
	assert.Equal(t, mockLogger, gateway.log)
}

func TestGetCardPrice_Success_NonFoil(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	responseBody := `{
		"prices": {
			"usd": "10.50",
			"usd_foil": "25.00"
		}
	}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.NoError(t, err)
	assert.Equal(t, 10.50, price)
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_Success_Foil(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            true,
	}

	responseBody := `{
		"prices": {
			"usd": "10.50",
			"usd_foil": "25.00"
		}
	}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.NoError(t, err)
	assert.Equal(t, 25.00, price)
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_RequestError(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(nil, fmt.Errorf("request error"))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "card gateway failed to get card")
	mockWeb.AssertExpectations(t)
}

func TestGetCardPrice_DoRequestError(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(nil, fmt.Errorf("do error"))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "card gateway failed to get response")
	mockWeb.AssertExpectations(t)
}

func TestGetCardPrice_CardNotFound(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "999",
		Foil:            false,
	}

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/999", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusNotFound)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader("")))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.IsType(t, ErrCardNotFound{}, err)
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_NonOKStatus(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	responseBody := "internal server error"

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusInternalServerError)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "card gateway failed to get card: http status 500")
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_InvalidJSON(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	responseBody := `{"invalid": json}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "card gateway failed to unmarshal body")
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_PriceIsZero_NonFoil(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	responseBody := `{
		"prices": {
			"usd": null,
			"usd_foil": "25.00"
		}
	}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.IsType(t, ErrPriceIsZero{}, err)
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_PriceIsZero_Foil(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            true,
	}

	responseBody := `{
		"prices": {
			"usd": "10.50",
			"usd_foil": null
		}
	}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.IsType(t, ErrPriceIsZero{}, err)
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}

func TestGetCardPrice_InvalidFloatParsing(t *testing.T) {
	mockWeb := mocks.NewHTTPMock()
	mockLogger := mocks.NewLogMock()
	mockRequest := mocks.NewRequestMock()
	mockResponse := mocks.NewResponseMock()

	gateway := New(mockWeb, mockLogger)

	card := domain.Cards{
		SetName:         "alpha",
		CollectorNumber: "161",
		Foil:            false,
	}

	responseBody := `{
		"prices": {
			"usd": "not-a-number",
			"usd_foil": "25.00"
		}
	}`

	mockWeb.On("NewRequestWithContext", mock.Anything, "GET", "https://api.scryfall.com/cards/alpha/161", mock.Anything).Return(mockRequest, nil)
	mockWeb.On("Do", mockRequest).Return(mockResponse, nil)
	mockResponse.On("StatusCode").Return(http.StatusOK)
	mockResponse.On("Body").Return(io.NopCloser(strings.NewReader(responseBody)))

	price, err := gateway.GetCardPrice(context.Background(), card)

	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "card gateway failed to parse float for usd")
	mockWeb.AssertExpectations(t)
	mockResponse.AssertExpectations(t)
}
