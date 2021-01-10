package favicon

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/db"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func downloadFavicon(domain *db.Domain, hostname string) *db.Domain {
	if domain == nil {
		domain = &db.Domain{
			Hostname: hostname,
		}
	}

	// lookup index.html (or whatever html is served) first
	iconURL, err := lookupIndex(hostname)

	if err != nil {
		logbuch.Debug("Error loading index to look up favicon", logbuch.Fields{
			"err":      err,
			"hostname": hostname,
		})

		// fallback to favicon.ico
		iconURL = filepath.Join(hostname, "favicon.ico")
	}

	file, err := downloadIcon(iconURL, hostname)

	if err != nil {
		logbuch.Debug("Error downloading favicon", logbuch.Fields{
			"err":      err,
			"hostname": hostname,
		})
	}

	domain.Filename = sql.NullString{String: file, Valid: file != ""}
	db.SaveDomain(nil, domain)
	return domain
}

func lookupIndex(hostname string) (string, error) {
	resp, err := http.Get("http://" + hostname)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)

	if err != nil {
		return "", err
	}

	head := findHeadNode(doc)

	if head == nil {
		return "", errors.New("error finding head node")
	}

	icons := findLinkIconNodes(head)

	if len(icons) == 0 {
		return "", errors.New("error finding link nodes")
	}

	// TODO select by preferred size
	for _, icon := range icons {
		href := getHref(&icon)

		if href != "" {
			return href, nil
		}
	}

	return "", errors.New("error finding href attribute")
}

func findHeadNode(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "head" {
		return node
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if n := findHeadNode(c); n != nil {
				return n
			}
		}
	}

	return nil
}

func findLinkIconNodes(headNode *html.Node) []html.Node {
	nodes := make([]html.Node, 0)

	for c := headNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "link" && hasRelAttribute(c) {
			nodes = append(nodes, *c)
		}
	}

	return nodes
}

func hasRelAttribute(node *html.Node) bool {
	for _, attr := range node.Attr {
		if strings.ToLower(attr.Key) == "rel" {
			val := strings.ToLower(attr.Val)

			if val == "icon" || val == "shortcut icon" || val == "apple-touch-icon" {
				return true
			}
		}
	}

	return false
}

func getHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if strings.ToLower(attr.Key) == "href" {
			return attr.Val
		}
	}

	return ""
}

func downloadIcon(rawurl, hostname string) (string, error) {
	resp, err := http.Get(getFaviconURL(rawurl, hostname))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	file, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Join(filesDir, hostname), 0744); err != nil {
		return "", err
	}

	filename := getFaviconFilename(rawurl)
	p := filepath.Join(filesDir, hostname, filename)

	if err := ioutil.WriteFile(p, file, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func getFaviconURL(rawurl, hostname string) string {
	u, err := url.Parse(rawurl)

	if err != nil {
		return ""
	}

	if u.Host == "" {
		u.Host = hostname
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return u.String()
}

func getFaviconFilename(rawurl string) string {
	queryParams := strings.Index(rawurl, "?")

	if queryParams > -1 {
		rawurl = rawurl[:queryParams]
	}

	return path.Base(rawurl)
}
