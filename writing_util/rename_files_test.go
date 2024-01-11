package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\2dd23015-d80e-4fe7-8535-f62b35b944e7_Export-4e406b8b-26c5-4d50-b179-d0a841dffb0f\6a1539bc4f92415fbc69d6d93091bdd6`, ".png")
	assert.Nil(t, err)
}
