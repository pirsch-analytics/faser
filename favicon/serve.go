package favicon

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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
	sizeParam := strings.TrimSpace(r.URL.Query().Get("size"))
	var size int
	var err error

	if sizeParam != "" {
		size, err = strconv.Atoi(sizeParam)
	}

	if hostname == "" || err != nil {
		sendBadRequest(w, err != nil)
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

	filename, size := selectFilenameForSize(domain.Filename.String, size)
	filePath := filepath.Join(filesDir, hostname, filename)
	_, err = os.Stat(filePath)

	if size != 0 && os.IsNotExist(err) {
		if err := scale(hostname, filename, size); err != nil {
			serveDefaultFavicon(w, r)
			return
		}
	}

	http.ServeFile(w, r, filePath)
}

func sendBadRequest(w http.ResponseWriter, sizeErr bool) {
	w.WriteHeader(http.StatusBadRequest)
	resp := struct {
		URL  string `json:"url"`
		Size string `json:"size,omitempty"`
	}{
		URL: "provide a valid URL or hostname",
	}

	if sizeErr {
		resp.Size = "provide a number greater or equal to 0"
	}

	respBody, err := json.Marshal(&resp)

	if err != nil {
		logbuch.Error("Error encoding error response", logbuch.Fields{"err": err})
		return
	}

	if _, err := w.Write(respBody); err != nil && err != syscall.EPIPE {
		logbuch.Warn("Error sending error response", logbuch.Fields{"err": err})
	}
}

func serveDefaultFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, server.Config().DefaultFavicon)
}
