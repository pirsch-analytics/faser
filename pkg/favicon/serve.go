package favicon

import (
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/pkg/server"
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
	setResponseHeaders(w, hostname)

	if hostname == "" {
		serveDefaultFavicon(w, r, fallback)
		return
	}

	var size int

	if sizeParam != "" {
		size, _ = strconv.Atoi(sizeParam)
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
		logbuch.Debug("Serving fallback icon", logbuch.Fields{"path": filepath.Join(server.Config().Cache.DefaultFaviconDir, server.Config().Cache.DefaultFavicon)})
		http.ServeFile(w, r, filepath.Join(server.Config().Cache.DefaultFaviconDir, server.Config().Cache.DefaultFavicon))
	} else {
		logbuch.Debug("Serving fallback icon", logbuch.Fields{"path": filepath.Join(server.Config().Cache.DefaultFaviconDir, filepath.Base(fallback))})
		http.ServeFile(w, r, filepath.Join(server.Config().Cache.DefaultFaviconDir, filepath.Base(fallback)))
	}
}
