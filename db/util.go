package db

import (
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

// Commit commits the given transaction and logs errors.
func Commit(tx *sqlx.Tx) {
	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction", logbuch.Fields{"err": err})
	}
}

// Rollback rolls back the given transaction and logs errors.
func Rollback(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		logbuch.Error("Error rolling back transaction", logbuch.Fields{"err": err})
	}
}

// CloseRows closes the given rows and logs errors.
func CloseRows(rows *sqlx.Rows) {
	if err := rows.Close(); err != nil {
		logbuch.Error("Error closing rows", logbuch.Fields{"err": err})
	}
}
