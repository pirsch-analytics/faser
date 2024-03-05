package favicon

import (
	"github.com/pirsch-analytics/faser/pkg/server"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	server.LoadConfig()
	server.Config().Cache.DefaultFaviconDir = "../../static"
	code := m.Run()
	os.Exit(code)
}
