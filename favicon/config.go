package favicon

import (
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/server"
)

var cache domainCache

// Init initializes the favicon cache.
func Init() {
	cache = domainCache{
		entries:    make(map[string]*db.Domain),
		maxAge:     server.Config().CacheMaxAge,
		maxEntries: server.Config().CacheMaxEntries,
	}
	cache.clear()
}
