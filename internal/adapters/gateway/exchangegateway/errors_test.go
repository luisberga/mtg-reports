package exchangegateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrFailedToGetExchangeRequest_Error(t *testing.T) {
	err := ErrFailedToGetExchangeRequest{}
	expected := "exchange gateway failed to get exchange request"

	assert.Equal(t, expected, err.Error())
}

func TestErrFailedToGetExchangeRequest_Interface(t *testing.T) {
	var err error = ErrFailedToGetExchangeRequest{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrExchangeRequestMultipleValues_Error(t *testing.T) {
	err := ErrExchangeRequestMultipleValues{}
	expected := "exchange gateway fund multiple values in exchange request"

	assert.Equal(t, expected, err.Error())
}

func TestErrExchangeRequestMultipleValues_Interface(t *testing.T) {
	var err error = ErrExchangeRequestMultipleValues{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrExchangeRequestNillValue_Error(t *testing.T) {
	err := ErrExchangeRequestNillValue{}
	expected := "exchange gateway fund nill value in exchange request"

	assert.Equal(t, expected, err.Error())
}

func TestErrExchangeRequestNillValue_Interface(t *testing.T) {
	var err error = ErrExchangeRequestNillValue{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrorTypes_AreDistinct(t *testing.T) {
	err1 := ErrFailedToGetExchangeRequest{}
	err2 := ErrExchangeRequestMultipleValues{}
	err3 := ErrExchangeRequestNillValue{}

	assert.IsType(t, ErrFailedToGetExchangeRequest{}, err1)
	assert.IsType(t, ErrExchangeRequestMultipleValues{}, err2)
	assert.IsType(t, ErrExchangeRequestNillValue{}, err3)

	// Ensure they are different types
	assert.NotEqual(t, err1, err2)
	assert.NotEqual(t, err1, err3)
	assert.NotEqual(t, err2, err3)
}
