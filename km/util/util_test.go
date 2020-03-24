package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	v, err := Load("data://YmFuYW5h")
	assert.NoError(t, err)
	assert.Equal(t, []byte("banana"), v)

	v, err = Load("file://testdata/load_data.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte("sasquatch"), v)
}