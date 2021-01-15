package favicon

import (
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/pirsch-analytics/faser/server"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func downloadFavicon(hostname string) string {
	// lookup index.html (or whatever html is served) first
	iconURL, err := lookupIndex(hostname)

	if err != nil {
		logbuch.Debug("Error loading index to look up favicon", logbuch.Fields{
			"err":      err,
			"hostname": hostname,
		})
	}

	// fallback to favicon.ico
	if iconURL == "" {
		iconURL = "favicon.ico"
	}

	filename, err := downloadIcon(iconURL, hostname)

	if err != nil {
		logbuch.Debug("Error downloading favicon", logbuch.Fields{
			"err":      err,
			"hostname": hostname,
		})
	}

	return filename
}

func lookupIndex(hostname string) (string, error) {
	resp, err := http.Get("http://" + hostname)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

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

	href := getHref(findLargestIconNode(icons))

	if href != "" {
		return href, nil
	}

	return "", errors.New("error finding href attribute")
}

func findHeadNode(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "head" {
		return node
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if n := findHeadNode(c); n != nil {
			return n
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

func findLargestIconNode(nodes []html.Node) *html.Node {
	width := 0
	var largest html.Node

	for _, n := range nodes {
		w := getSizesWidth(&n)

		if w >= width {
			width = w
			largest = n
		}
	}

	return &largest
}

func getSizesWidth(node *html.Node) int {
	for _, attr := range node.Attr {
		if strings.ToLower(attr.Key) == "sizes" && attr.Val != "" {
			val := strings.ToLower(attr.Val)
			sizes := strings.Split(val, " ")

			if len(sizes) > 0 {
				dimension := strings.Split(sizes[len(sizes)-1], "x")

				if len(dimension) == 2 {
					width, err := strconv.Atoi(dimension[1])

					if err != nil {
						return 0
					}

					return width
				}
			}

			return 0
		}
	}

	return 0
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

	dir := server.Config().Cache.Dir

	if err := os.MkdirAll(filepath.Join(dir, hostname), 0744); err != nil {
		return "", err
	}

	filename := getFaviconFilename(rawurl)
	ext := path.Ext(filename)
	filename = "favicon" + ext
	p := filepath.Join(dir, hostname, filename)

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
