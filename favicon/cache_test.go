package favicon

import (
	"database/sql"
	"github.com/pirsch-analytics/faser/db"
	"testing"
	"time"
)

func TestCacheNotFound(t *testing.T) {
	Init()
	clearDB(t)
	cache := domainCache{
		entries:    make(map[string]*db.Domain),
		maxAge:     100,
		maxEntries: 100,
	}
	cache.clear()

	if cache.get("example.com") != nil {
		t.Fatal("Domain must not have been found")
	}

	if len(cache.entries) != 0 {
		t.Fatalf("Cache must have no entries, but was: %v", len(cache.entries))
	}
}

func TestCacheFound(t *testing.T) {
	Init()
	clearDB(t)
	cache := domainCache{
		entries:    make(map[string]*db.Domain),
		maxAge:     1,
		maxEntries: 100,
	}
	cache.clear()
	db.SaveDomain(nil, &db.Domain{
		Hostname: "example.com",
		Filename: sql.NullString{String: "favicon.png", Valid: true},
	})

	if cache.get("example.com") == nil {
		t.Fatal("Domain must have been found")
	}

	if len(cache.entries) != 1 {
		t.Fatalf("Cache must have one entry, but was: %v", len(cache.entries))
	}

	if cache.get("not-found.com") != nil {
		t.Fatal("Domain must not have been found")
	}

	if len(cache.entries) != 1 {
		t.Fatalf("Cache must have one entry, but was: %v", len(cache.entries))
	}
}

func TestCacheMaxAge(t *testing.T) {
	Init()
	clearDB(t)
	cache := domainCache{
		entries:    make(map[string]*db.Domain),
		maxAge:     1,
		maxEntries: 100,
	}
	cache.clear()
	db.SaveDomain(nil, &db.Domain{
		Hostname: "example.com",
		Filename: sql.NullString{String: "favicon.png", Valid: true},
	})
	db.SaveDomain(nil, &db.Domain{
		Hostname: "example2.com",
		Filename: sql.NullString{String: "favicon.png", Valid: true},
	})

	if cache.get("example.com") == nil {
		t.Fatal("Domain must have been found")
	}

	if cache.get("example2.com") == nil {
		t.Fatal("Domain must have been found")
	}

	if len(cache.entries) != 2 {
		t.Fatalf("Cache must have two entries, but was: %v", len(cache.entries))
	}

	time.Sleep(1 * time.Second)

	if cache.get("example.com") != nil {
		t.Fatal("Domain must not have been found")
	}

	if len(cache.entries) != 1 {
		t.Fatalf("Cache must have one entry, but was: %v", len(cache.entries))
	}
}

func TestCacheMaxEntries(t *testing.T) {
	Init()
	clearDB(t)
	cache := domainCache{
		entries:    make(map[string]*db.Domain),
		maxAge:     100,
		maxEntries: 1,
	}
	cache.clear()
	db.SaveDomain(nil, &db.Domain{
		Hostname: "example.com",
		Filename: sql.NullString{String: "favicon.png", Valid: true},
	})
	db.SaveDomain(nil, &db.Domain{
		Hostname: "example2.com",
		Filename: sql.NullString{String: "favicon.png", Valid: true},
	})

	if cache.get("example.com") == nil {
		t.Fatal("Domain must have been found")
	}

	if cache.get("example2.com") == nil {
		t.Fatal("Domain must have been found")
	}

	if len(cache.entries) != 1 {
		t.Fatalf("Cache must have one entry, but was: %v", len(cache.entries))
	}
}

func clearDB(t *testing.T) {
	if _, err := db.Exec(nil, `DELETE FROM "domain"`); err != nil {
		t.Fatal(err)
	}
}
