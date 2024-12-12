package lruk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	l, err := NewLRUK(10, 3)
	assert.Nil(t, err)

	value, ok := l.Get(String("key1"))
	assert.Falsef(t, ok, "value should not be found")
	assert.Nil(t, value)
}
