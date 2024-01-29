package writing_util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeSpaceFilesInFolder(t *testing.T) {
	err := removeSpaceFilesInFolder(`C:\Users\caijiatao\Downloads\bec95a90-d792-4de9-8932-e8133af4b45a_Export-84c8354d-16f8-498b-bc6a-2c9c2e416f96\8a369eac176c437fafd3592e8e0dd2fa`, ".png")
	assert.Nil(t, err)
}
