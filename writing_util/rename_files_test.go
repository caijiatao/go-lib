package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\da2937ab-94a5-4190-8192-637145ab608f_Export-4ccef35a-ad18-4ef6-af75-f8c62d1de63c\a7cec09e91b34db4b501776888c99b8a`, ".png")
	assert.Nil(t, err)
}
