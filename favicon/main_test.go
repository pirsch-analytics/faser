package favicon

import (
	"github.com/pirsch-analytics/faser/server"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server.LoadConfig()
	code := m.Run()
	os.Exit(code)
}
