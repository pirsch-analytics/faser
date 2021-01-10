package favicon

import (
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

var sizes = []int{
	16,
	32,
	64,
	96,
	128,
	196,
}

// ServeFavicon looks up a favicon and serves the file in the desired dimensions if possible or the default icon otherwise.
func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	hostname := getHostname(r.URL.Query().Get("url"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

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

	filename := getFilenameForSize(domain.Filename.String, size)
	http.ServeFile(w, r, filepath.Join(filesDir, hostname, filename))
}

func serveDefaultFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, server.Config().DefaultFavicon)
}

func getFilenameForSize(filename string, size int) string {
	if size <= 0 {
		return filename
	}

	for _, s := range sizes {
		if size <= s {
			size = s
			break
		}
	}

	if size > sizes[len(sizes)-1] {
		size = sizes[len(sizes)-1]
	}

	ext := path.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	return filenameWithoutExt + "-" + strconv.Itoa(size) + ext
}
