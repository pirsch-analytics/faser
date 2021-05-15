package favicon

import (
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// ServeFavicon looks up a favicon and serves the file in the desired dimensions if possible or the default icon otherwise.
func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname(r.URL.Query().Get("url"))
	sizeParam := strings.TrimSpace(r.URL.Query().Get("size"))
	fallback := strings.TrimSpace(r.URL.Query().Get("fallback"))
	var size int
	var err error

	if sizeParam != "" {
		size, err = strconv.Atoi(sizeParam)
	}

	setResponseHeaders(w, hostname)

	if hostname == "" || err != nil {
		serveDefaultFavicon(w, r, fallback)
		return
	}

	filename := faviconCache.find(hostname, size)

	if filename == "" {
		serveDefaultFavicon(w, r, fallback)
		return
	}

	http.ServeFile(w, r, filepath.Join(server.Config().Cache.Dir, hostname, filename))
}

func setResponseHeaders(w http.ResponseWriter, etag string) {
	w.Header().Add("Permissions-Policy", "interest-cohort=()")
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Header().Add("ETag", etag)
}

func serveDefaultFavicon(w http.ResponseWriter, r *http.Request, fallback string) {
	if fallback == "" {
		http.ServeFile(w, r, server.Config().Cache.DefaultFavicon)
	} else {
		http.ServeFile(w, r, filepath.Join(server.Config().Cache.DefaultFaviconDir, filepath.Base(fallback)))
	}
}
