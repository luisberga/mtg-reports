package reporthandler

import (
	"context"
	"fmt"
	"testing"

	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	mockReportService := mocks.NewReportServiceMock()
	mockLogger := mocks.NewLogMock()

	handler := New(mockReportService, mockLogger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockReportService, handler.ReportService)
	assert.Equal(t, mockLogger, handler.log)
}

func TestProcessAndSend_Success(t *testing.T) {
	mockReportService := mocks.NewReportServiceMock()
	mockLogger := mocks.NewLogMock()

	handler := New(mockReportService, mockLogger)

	mockLogger.On("Info", mock.Anything).Twice()
	mockReportService.On("ProcessAndSend", mock.Anything).Return(nil)

	err := handler.ProcessAndSend(context.Background())

	assert.NoError(t, err)
	mockReportService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestProcessAndSend_ServiceError(t *testing.T) {
	mockReportService := mocks.NewReportServiceMock()
	mockLogger := mocks.NewLogMock()

	handler := New(mockReportService, mockLogger)

	expectedError := fmt.Errorf("service error")

	mockLogger.On("Info", mock.Anything).Once()
	mockReportService.On("ProcessAndSend", mock.Anything).Return(expectedError)

	err := handler.ProcessAndSend(context.Background())

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockReportService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
