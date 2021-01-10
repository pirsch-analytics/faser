package favicon

import (
	"github.com/emvi/logbuch"
	"net/url"
)

func getHostname(rawurl string) string {
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
