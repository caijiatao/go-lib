package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	db, err := initDB()
	assert.Nil(t, err)
	err = queryExample(db)
	assert.Nil(t, err)
}
