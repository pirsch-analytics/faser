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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	domain := db.GetDomain(hostname)
	maxAge := time.Now().UTC().Add(-time.Duration(server.Get().Cache) * time.Second)

	if domain == nil || domain.ModTime.Before(maxAge) {
		domain = downloadFavicon(domain, hostname)
	}

	if !domain.Filename.Valid {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath.Join(filesDir, hostname, domain.Filename.String))
}
