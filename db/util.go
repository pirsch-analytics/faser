package db

import (
	"database/sql"
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

// Exec executes a single SQL statement in given transaction or creates a new one if nil.
// This function should not be used for production code and is intended to be used for tests only.
func Exec(tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	if tx == nil {
		tx, _ = pool.Beginx()
		defer Commit(tx)
	}

	result, err := tx.Exec(query, args...)

	if err != nil {
		logbuch.Error("Error executing sql statement", logbuch.Fields{"err": err})
		return result, err
	}

	return result, nil
}

// Tx creates a new transaction or panics on failure.
func Tx() *sqlx.Tx {
	tx, err := pool.Beginx()

	if err != nil {
		logbuch.Fatal("Error creating transaction", logbuch.Fields{"err": err})
	}

	return tx
}
