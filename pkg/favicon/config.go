package favicon

import (
	"github.com/pirsch-analytics/faser/pkg/server"
	"time"
)

var faviconCache *cache

// Init initializes the favicon cache.
func Init() {
	faviconCache = newCache(server.Config().Cache.MaxEntries, time.Duration(server.Config().Cache.MaxAge)*time.Second)
}
