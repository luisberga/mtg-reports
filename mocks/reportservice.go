package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ReportServiceMock struct {
	mock.Mock
}

func NewReportServiceMock() *ReportServiceMock {
	return &ReportServiceMock{}
}

func (m *ReportServiceMock) ProcessAndSend(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
