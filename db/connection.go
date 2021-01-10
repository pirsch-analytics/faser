package db

import (
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pirsch-analytics/faser/server"
)

const (
	connectTimeout            = "60"
	defaultMaxOpenConnections = 1
	connectionString          = `host=%s
		port=%s
		user=%s
		password=%s
		dbname=%s
		sslmode=%s
		sslcert=%s
		sslkey=%s
		sslrootcert=%s
		connectTimeout=%s
		timezone=%s`
)

type connection struct {
	sqlx.DB
}

func newConnection() *connection {
	logbuch.Info("Connecting to database...", logbuch.Fields{"ssl_mode": server.Config().DB.SSLMode})
	db, err := sqlx.Connect("postgres", postgresConnection())

	if err != nil {
		logbuch.Fatal("Error connecting to database", logbuch.Fields{"err": err})
		return nil
	}

	if err := db.Ping(); err != nil {
		logbuch.Fatal("Error pinging database", logbuch.Fields{"err": err})
		return nil
	}

	logbuch.Info("Connected")
	setMaxOpenConnections(db, server.Config().DB.MaxOpenConnections)
	return &connection{DB: *db}
}

func (connection *connection) disconnect() {
	logbuch.Info("Disconnecting from database...")

	if err := connection.Close(); err != nil {
		logbuch.Warn("Error when closing database connection", logbuch.Fields{"err": err})
	}

	logbuch.Info("Disconnected")
}

func postgresConnection() string {
	data := server.Config().DB
	return fmt.Sprintf(connectionString, data.Host, data.Port, data.User, data.Password, data.Schema,
		data.SSLMode, data.SSLCert, data.SSLKey, data.SSLRootCert, connectTimeout, "UTC")
}

func setMaxOpenConnections(db *sqlx.DB, maxOpenConnections int) {
	if maxOpenConnections == 0 {
		maxOpenConnections = defaultMaxOpenConnections
	}

	logbuch.Info("Setting max open connections to database", logbuch.Fields{"max_open_connections": maxOpenConnections})
	db.SetMaxOpenConns(maxOpenConnections)
}
