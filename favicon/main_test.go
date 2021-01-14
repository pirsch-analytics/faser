package favicon

import (
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/server"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server.LoadConfig()
	server.Config().DB.Host = "localhost"
	server.Config().DB.Port = "5432"
	server.Config().DB.Schema = "faser-test"
	server.Config().DB.User = "postgres"
	server.Config().DB.Password = "postgres"
	server.Config().DB.SSLMode = "disable"
	server.Config().DB.MigrationDir = "../schema"
	db.Migrate()
	db.Connect()
	code := m.Run()
	db.Disconnect()
	os.Exit(code)
}
