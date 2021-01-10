package favicon

import (
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"path/filepath"
	"time"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname(r.URL.Query().Get("url"))

	if hostname == "" {
		serveDefaultFavicon(w, r)
		return
	}

	domain := db.GetDomain(hostname)
	maxAge := time.Now().UTC().Add(-time.Duration(server.Config().Cache) * time.Second)

	if domain == nil || domain.ModTime.Before(maxAge) {
		domain = downloadFavicon(domain, hostname)
	}

	if !domain.Filename.Valid {
		serveDefaultFavicon(w, r)
		return
	}

	http.ServeFile(w, r, filepath.Join(filesDir, hostname, domain.Filename.String))
}

func serveDefaultFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, server.Config().DefaultFavicon)
}
