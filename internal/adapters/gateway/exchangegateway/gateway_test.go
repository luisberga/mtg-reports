package exchangegateway

import (
	"context"
	"errors"
	"io"
	"mtg-report/mocks"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	logMock := mocks.NewLogMock()
	url := "https://api.exchangerate-api.com/v4/latest/USD"

	gateway := New(webMock, url, logMock)

	assert.NotNil(t, gateway)
	assert.Equal(t, url, gateway.url)
}

func TestExchangeGateway_GetUSD_Success(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	respMock := mocks.NewResponseMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	responseBody := `{"conversion_rates":{"BRL":5.25}}`
	bodyReader := io.NopCloser(strings.NewReader(responseBody))

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(respMock, nil)
	respMock.On("StatusCode").Return(http.StatusOK)
	respMock.On("Body").Return(bodyReader)

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 5.25, got)

	webMock.AssertExpectations(t)
	respMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_RequestCreationError(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	logMock := mocks.NewLogMock()

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).
		Return(nil, errors.New("request creation failed"))

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.Contains(t, err.Error(), "exchange gateway failed to create request")

	webMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_RequestExecutionError(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(nil, errors.New("request execution failed"))

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.Contains(t, err.Error(), "exchange gateway failed to get response")

	webMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_NonOKStatusCode(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	respMock := mocks.NewResponseMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	bodyReader := io.NopCloser(strings.NewReader(""))

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(respMock, nil)
	respMock.On("StatusCode").Return(http.StatusInternalServerError)
	respMock.On("Body").Return(bodyReader)

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.IsType(t, ErrFailedToGetExchangeRequest{}, err)

	webMock.AssertExpectations(t)
	respMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_InvalidJSON(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	respMock := mocks.NewResponseMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	responseBody := `{"invalid json"`
	bodyReader := io.NopCloser(strings.NewReader(responseBody))

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(respMock, nil)
	respMock.On("StatusCode").Return(http.StatusOK)
	respMock.On("Body").Return(bodyReader)

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.Contains(t, err.Error(), "exchange gateway failed to unmarshal body")

	webMock.AssertExpectations(t)
	respMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_NilBRLRate(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	respMock := mocks.NewResponseMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	responseBody := `{"conversion_rates":{"BRL":null}}`
	bodyReader := io.NopCloser(strings.NewReader(responseBody))

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(respMock, nil)
	respMock.On("StatusCode").Return(http.StatusOK)
	respMock.On("Body").Return(bodyReader)

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.IsType(t, ErrExchangeRequestNillValue{}, err)

	webMock.AssertExpectations(t)
	respMock.AssertExpectations(t)
}

func TestExchangeGateway_GetUSD_MissingBRLRate(t *testing.T) {
	webMock := mocks.NewHTTPMock()
	respMock := mocks.NewResponseMock()
	reqMock := mocks.NewRequestMock()
	logMock := mocks.NewLogMock()

	responseBody := `{"conversion_rates":{"USD":1.0}}`
	bodyReader := io.NopCloser(strings.NewReader(responseBody))

	webMock.On("NewRequestWithContext", mock.Anything, "GET", "https://api.test.com", nil).Return(reqMock, nil)
	webMock.On("Do", reqMock).Return(respMock, nil)
	respMock.On("StatusCode").Return(http.StatusOK)
	respMock.On("Body").Return(bodyReader)

	gateway := New(webMock, "https://api.test.com", logMock)
	got, err := gateway.GetUSD(context.Background())

	assert.Error(t, err)
	assert.Equal(t, float64(0), got)
	assert.IsType(t, ErrExchangeRequestNillValue{}, err)

	webMock.AssertExpectations(t)
	respMock.AssertExpectations(t)
}
