package favicon

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

	if hostname == "" || err != nil {
		sendBadRequest(w, err != nil)
		return
	}

	filename := faviconCache.find(hostname, size)

	if filename == "" {
		serveDefaultFavicon(w, r, fallback)
		return
	}

	http.ServeFile(w, r, filepath.Join(server.Config().Cache.Dir, hostname, filename))
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

func serveDefaultFavicon(w http.ResponseWriter, r *http.Request, fallback string) {
	if fallback == "" {
		http.ServeFile(w, r, server.Config().Cache.DefaultFavicon)
	} else {
		http.ServeFile(w, r, filepath.Join(server.Config().Cache.DefaultFaviconDir, filepath.Base(fallback)))
	}
}
