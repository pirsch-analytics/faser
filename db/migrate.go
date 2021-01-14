package db

import (
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pirsch-analytics/faser/server"

	// migration database driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	// migrate from source files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrateConnectionString = `postgres://%s:%s/%s?user=%s&password=%s&sslmode=%s&sslcert=%s&sslkey=%s&sslrootcert=%s&connect_timeout=60`
	migrationDir            = "./schema"
)

// Migrate runs database migrations scripts and panics in case of an error.
func Migrate() {
	logbuch.Info("Migrating database schema...", logbuch.Fields{"dir": migrationDir})
	m, err := migrate.New(
		"file://"+migrationDir,
		postgresConnectionString())

	if err != nil {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"err": err})
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"err": err})
		return
	}

	if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"source_err": sourceErr, "db_err": dbErr})
	}

	logbuch.Info("Successfully migrated database schema")
}

func postgresConnectionString() string {
	cfg := server.Config().DB
	return fmt.Sprintf(migrateConnectionString, cfg.Host, cfg.Port, cfg.Schema, cfg.User, cfg.Password,
		cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)
}
