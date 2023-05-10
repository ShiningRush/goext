package errx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchErrors(t *testing.T) {
	err := &BatchErrors{}
	err.Append(errors.New("err1"))
	err.Append(errors.New("err2"))

	assert.Equal(t, 2, err.Len())
	assert.True(t, err.HasError())
	assert.Equal(t, "\nerr[0]: err1\nerr[1]: err2", err.Error())
}
