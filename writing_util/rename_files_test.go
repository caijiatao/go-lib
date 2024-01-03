package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\c71233a1-3363-439e-b90d-be89a18dfb61_Export-d1f4e9fd-164e-47e1-af76-71dfd3edbe4b\软件工程 b9abd1115c8c4628b7a13202a4f89823`, ".png")
	assert.Nil(t, err)
}
