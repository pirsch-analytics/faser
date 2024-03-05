package favicon

import (
	"net/url"
)

func getHostname(rawURL string) string {
	u, err := url.ParseRequestURI(rawURL)

	if err != nil || u.Host == "" {
		rawURL = "http://" + rawURL
		u, err = url.ParseRequestURI(rawURL)

		if err != nil {
			return ""
		}
	}

	return u.Hostname()
}
