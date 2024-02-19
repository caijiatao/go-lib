package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\bc100d51-f259-4c3d-9cf1-662eccc1e4a6_Export-31ee2d10-bf51-4a8b-9e42-d99ad94de840\9b49f76db38d43f3a329b63c1d4d51d9`, ".png")
	assert.Nil(t, err)
}
