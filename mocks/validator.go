package mocks

import (
	"mtg-report/internal/core/dtos"

	"github.com/stretchr/testify/mock"
)

type ValidateMock struct {
	mock.Mock
}

func NewValidateMock() *ValidateMock {
	return &ValidateMock{}
}

func (v *ValidateMock) Card(card dtos.RequestInsertCard) error {
	args := v.Called(card)
	return args.Error(0)
}

func (v *ValidateMock) CardID(parts []string) (string, error) {
	args := v.Called(parts)
	return args.String(0), args.Error(1)
}

func (v *ValidateMock) Filters(setName, name, collector_number string) map[string]string {
	args := v.Called(setName, name, collector_number)
	return args.Get(0).(map[string]string)
}

func (v *ValidateMock) CardName(card dtos.RequestUpdateCard) error {
	args := v.Called(card)
	return args.Error(0)
}

func (v *ValidateMock) Pagination(pageStr, limitStr string) (int, int, error) {
	args := v.Called(pageStr, limitStr)
	return args.Int(0), args.Int(1), args.Error(2)
}
