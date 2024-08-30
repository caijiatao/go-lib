package good_t

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseInt(t *testing.T) {
	result := ParseInt("5")
	assert.Equal(t, 5, result)
}
