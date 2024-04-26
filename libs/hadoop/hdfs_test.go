package hadoop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	err := WriteFile()
	assert.Nil(t, err)
	ReadFile()

}
