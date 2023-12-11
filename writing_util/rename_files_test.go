package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\9248ba84-818a-42dc-9df3-7b0d4b62ac7c_Export-4410964d-3ad7-4949-afa0-082821a5d1d9\一文搞懂kubernetes 中的负载均衡 f44ed98fb5c944e1b29af6c2c3e4f617`, ".png")
	assert.Nil(t, err)
}
