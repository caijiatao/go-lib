package orm

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

func IsDuplicateErr(err error) bool {
	err = errors.Cause(err)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return true
	}
	return false
}
