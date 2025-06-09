package cardgateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrFailedToGetCardRequest_Error(t *testing.T) {
	err := ErrFailedToGetCardRequest{httpStatus: 500}
	expected := "failed to get card request: 500"
	assert.Equal(t, expected, err.Error())
}

func TestErrFailedToGetCardRequest_Interface(t *testing.T) {
	var err error = ErrFailedToGetCardRequest{httpStatus: 404}
	assert.NotNil(t, err)
}

func TestErrPriceIsZero_Error(t *testing.T) {
	err := ErrPriceIsZero{}
	expected := "price is zero - card could be foil or non-foil, check register"
	assert.Equal(t, expected, err.Error())
}

func TestErrPriceIsZero_Interface(t *testing.T) {
	var err error = ErrPriceIsZero{}
	assert.NotNil(t, err)
}

func TestErrCardNotFound_Error(t *testing.T) {
	err := ErrCardNotFound{}
	expected := "card not found"
	assert.Equal(t, expected, err.Error())
}

func TestErrCardNotFound_Interface(t *testing.T) {
	var err error = ErrCardNotFound{}
	assert.NotNil(t, err)
}

func TestErrorTypes_AreDistinct(t *testing.T) {
	err1 := ErrFailedToGetCardRequest{httpStatus: 500}
	err2 := ErrPriceIsZero{}
	err3 := ErrCardNotFound{}

	assert.NotEqual(t, err1.Error(), err2.Error())
	assert.NotEqual(t, err1.Error(), err3.Error())
	assert.NotEqual(t, err2.Error(), err3.Error())
}
