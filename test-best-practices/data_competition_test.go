package test_best_practices

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestIncrement(t *testing.T) {
	count := 100
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			increment()
			wg.Done()
		}()
	}
	assert.Equal(t, count, counter)
}
