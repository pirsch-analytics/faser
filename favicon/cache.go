package favicon

import (
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/server"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	filename string
	created  time.Time
}

type cache struct {
	favicons map[string]map[int]cacheEntry
	maxSize  int
	maxAge   time.Duration
	m        sync.RWMutex
}

func newCache(maxSize int, maxAge time.Duration) *cache {
	c := &cache{
		favicons: make(map[string]map[int]cacheEntry),
		maxSize:  maxSize,
		maxAge:   maxAge,
	}
	c.clear()
	return c
}

func (c *cache) find(hostname string, size int) string {
	if c.size() > c.maxSize {
		c.clear()
	}

	filename, found := c.get(hostname, size)

	log.Printf("%s -> %s", hostname, filename)

	if !found {
		filename = downloadFavicon(hostname)

		if filename != "" && size != 0 && scalableType(filepath.Ext(filename)) {
			filename = scale(hostname, filename, size)
		}

		c.set(hostname, size, filename)
	}

	return filename
}

// needs to return if the entry was found, because it might be empty (file couldn't be downloaded)
func (c *cache) get(hostname string, size int) (string, bool) {
	c.m.RLock()
	c.m.RUnlock()
	sizes, found := c.favicons[hostname]

	if !found {
		return "", false
	}

	favicon, found := sizes[size]

	if !found {
		return "", false
	}

	if favicon.created.Before(time.Now().Add(-c.maxAge)) {
		return "", false
	}

	return favicon.filename, true
}

func (c *cache) set(hostname string, size int, filename string) {
	c.m.Lock()
	defer c.m.Unlock()
	favicon, found := c.favicons[hostname]

	if !found {
		c.favicons[hostname] = map[int]cacheEntry{size: {filename, time.Now()}}
	} else {
		favicon[size] = cacheEntry{filename, time.Now()}
	}
}

func (c *cache) size() int {
	c.m.RLock()
	c.m.RUnlock()
	return len(c.favicons)
}

func (c *cache) clear() {
	c.m.Lock()
	defer c.m.Unlock()
	c.favicons = make(map[string]map[int]cacheEntry)

	if err := os.RemoveAll(server.Config().Cache.Dir); err != nil {
		logbuch.Fatal("Error deleting cache directory", logbuch.Fields{"err": err})
	}

	if err := os.MkdirAll(server.Config().Cache.Dir, 0744); err != nil {
		logbuch.Fatal("Error creating cache directory", logbuch.Fields{"err": err})
	}
}
