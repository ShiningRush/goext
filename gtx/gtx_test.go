package gtx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	assert.Equal(t, 0, Zero[int]())
	assert.Equal(t, "", Zero[string]())
	assert.Equal(t, struct{}{}, Zero[struct{}]())
}

func TestIsZero(t *testing.T) {
	testInt := 0
	assert.True(t, IsZero(testInt))
	testInt = 10
	assert.False(t, IsZero(testInt))

	testStr := ""
	assert.True(t, IsZero(testStr))
	testStr = "10"
	assert.False(t, IsZero(testStr))
}
