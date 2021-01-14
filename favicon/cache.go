package favicon

import (
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/db"
	"os"
	"sync"
	"time"
)

const (
	filesDir = "files"
)

type domainCache struct {
	entries    map[string]*db.Domain
	maxAge     int // seconds
	maxEntries int
	m          sync.RWMutex
}

func (cache *domainCache) get(hostname string) *db.Domain {
	cache.m.RLock()

	if len(cache.entries) >= cache.maxEntries {
		cache.m.RUnlock()
		cache.clear()
		cache.m.RLock()
	}

	domain, found := cache.entries[hostname]
	cache.m.RUnlock()

	if !found {
		domain = db.GetDomain(hostname)

		if domain != nil {
			cache.m.Lock()
			cache.entries[hostname] = domain
			cache.m.Unlock()
		}
	}

	if domain != nil {
		maxAge := time.Now().Add(-time.Duration(cache.maxAge) * time.Second)

		if domain.ModTime.Before(maxAge) {
			cache.m.Lock()
			delete(cache.entries, hostname)
			cache.m.Unlock()
			return nil
		}
	}

	return domain
}

func (cache *domainCache) clear() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.entries = make(map[string]*db.Domain)

	if err := os.RemoveAll(filesDir); err != nil {
		logbuch.Error("Error deleting files directory while clearing cache", logbuch.Fields{"err": err})
	}

	if err := os.MkdirAll(filesDir, 0744); err != nil {
		logbuch.Fatal("Error recreating files directory while clearing cache", logbuch.Fields{"err": err})
	}
}
