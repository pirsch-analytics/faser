package favicon

import (
	"github.com/emvi/logbuch"
	"os"
)

const (
	filesDir = "files"
)

// InitFileDir creates the file directory for favicons.
func InitFileDir() {
	if err := os.MkdirAll(filesDir, 0744); err != nil {
		logbuch.Fatal("Error creating file directory", logbuch.Fields{"err": err})
	}
}
