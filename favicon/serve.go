package favicon

import (
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := getDomain(vars["url"])

	if domain == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO
	// 1. look up domain (including max-age)
	// 2. download favicon if required
	// 3. serve
}

func getDomain(rawurl string) string {
	u, err := url.ParseRequestURI(rawurl)

	if err != nil || u.Host == "" {
		rawurl = "http://" + rawurl
		u, err = url.ParseRequestURI(rawurl)

		if err != nil {
			logbuch.Debug("Error parsing URL", logbuch.Fields{
				"err": err,
				"url": rawurl,
			})
			return ""
		}
	}

	return u.Hostname()
}
