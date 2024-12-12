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

	l.Add(String("key1"), String("1234"))
	value, ok = l.Get(String("key1"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)
}

func TestAdd(t *testing.T) {
	l, err := NewLRUK(3, 3)
	assert.Nil(t, err)

	l.Add(String("key1"), String("1234"))
	value, ok := l.Get(String("key1"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)

	l.Add(String("key2"), String("1234"))
	value, ok = l.Get(String("key2"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)

	l.Add(String("key3"), String("1234"))
	value, ok = l.Get(String("key3"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)

	l.Add(String("key4"), String("1234"))
	value, ok = l.Get(String("key4"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)

	value, ok = l.Get(String("key1"))
	assert.Falsef(t, ok, "value should not be found")
	assert.Nil(t, value)
}

func TestRemove(t *testing.T) {
	l, err := NewLRUK(3, 3)
	assert.Nil(t, err)

	l.Add(String("key1"), String("1234"))
	value, ok := l.Get(String("key1"))
	assert.Truef(t, ok, "value should be found")
	assert.Equal(t, String("1234"), value)

	l.Remove(String("key1"))
	value, ok = l.Get(String("key1"))
	assert.Falsef(t, ok, "value should not be found")
	assert.Nil(t, value)
}
