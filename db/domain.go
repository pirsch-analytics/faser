package db

import (
	"database/sql"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

// Domain represents a single domain entry in the database.
type Domain struct {
	BaseEntity

	Hostname string
	Filename sql.NullString
}

// GetDomain returns the domain for given hostname.
func GetDomain(hostname string) *Domain {
	query := `SELECT * FROM "domain" WHERE hostname = $1`
	entity := new(Domain)

	if err := pool.Get(entity, query, hostname); err != nil {
		logbuch.Debug("Domain not found", logbuch.Fields{
			"err":      err,
			"hostname": hostname,
		})
		return nil
	}

	return entity
}

// SaveDomain creates or updates a domain in the database.
func SaveDomain(tx *sqlx.Tx, entity *Domain) *Domain {
	pool.SaveEntity(tx, entity, `INSERT INTO "domain" (hostname, filename) VALUES (:hostname, :filename) RETURNING id`,
		`UPDATE "domain" SET hostname = :hostname, filename = :filename WHERE id = :id`)
	return entity
}

// DeleteDomain deletes all domains in the database.
func DeleteDomain(tx *sqlx.Tx) {
	if tx == nil {
		tx = Tx()
		defer Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "domain"`); err != nil {
		Rollback(tx)
		logbuch.Fatal("Error deleting domains", logbuch.Fields{"err": err})
	}
}
