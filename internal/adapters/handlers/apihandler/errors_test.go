package apihandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInternalErr_Error(t *testing.T) {
	err := ErrInternalErr{}
	expected := "internal error"

	assert.Equal(t, expected, err.Error())
}

func TestErrInternalErr_Interface(t *testing.T) {
	var err error = ErrInternalErr{}
	assert.NotNil(t, err)
	assert.Implements(t, (*error)(nil), err)
}

func TestErrInternalErr_Type(t *testing.T) {
	err := ErrInternalErr{}
	assert.IsType(t, ErrInternalErr{}, err)
}
