package cmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	cm := New(100)
	cm.Set("hello", "world")
	assert.Equal(t, cm.Get("hello"), "world", "Get != Set")
}
