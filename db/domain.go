package db

import (
	"database/sql"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Domain struct {
	BaseEntity

	Hostname string
	Filename sql.NullString
}

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

func SaveDomain(tx *sqlx.Tx, entity *Domain) *Domain {
	pool.SaveEntity(tx, entity, `INSERT INTO "domain" (hostname, filename) VALUES (:hostname, :filename) RETURNING id`,
		`UPDATE "domain" SET hostname = :hostname, filename = :filename WHERE id = :id`)
	return entity
}
