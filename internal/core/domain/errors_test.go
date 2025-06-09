package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrCardAlreadyExists_Error(t *testing.T) {
	err := ErrCardAlreadyExists{}
	expected := "card already exists"

	assert.Equal(t, expected, err.Error())
}

func TestErrCardAlreadyExists_Interface(t *testing.T) {
	var err error = ErrCardAlreadyExists{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrCardNotFound_Error(t *testing.T) {
	err := ErrCardNotFound{}
	expected := "card not found"

	assert.Equal(t, expected, err.Error())
}

func TestErrCardNotFound_Interface(t *testing.T) {
	var err error = ErrCardNotFound{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrCardsPriceNotFound_Error(t *testing.T) {
	err := ErrCardsPriceNotFound{}
	expected := "cards price not found"

	assert.Equal(t, expected, err.Error())
}

func TestErrCardsPriceNotFound_Interface(t *testing.T) {
	var err error = ErrCardsPriceNotFound{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrInvalidSetName_Error(t *testing.T) {
	err := ErrInvalidSetName{}
	expected := "invalid set name"

	assert.Equal(t, expected, err.Error())
}

func TestErrInvalidSetName_Interface(t *testing.T) {
	var err error = ErrInvalidSetName{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestDomainErrorTypes_AreDistinct(t *testing.T) {
	err1 := ErrCardAlreadyExists{}
	err2 := ErrCardNotFound{}
	err3 := ErrCardsPriceNotFound{}
	err4 := ErrInvalidSetName{}

	assert.IsType(t, ErrCardAlreadyExists{}, err1)
	assert.IsType(t, ErrCardNotFound{}, err2)
	assert.IsType(t, ErrCardsPriceNotFound{}, err3)
	assert.IsType(t, ErrInvalidSetName{}, err4)

	// Ensure they are different types
	assert.NotEqual(t, err1, err2)
	assert.NotEqual(t, err1, err3)
	assert.NotEqual(t, err1, err4)
	assert.NotEqual(t, err2, err3)
	assert.NotEqual(t, err2, err4)
	assert.NotEqual(t, err3, err4)
}
