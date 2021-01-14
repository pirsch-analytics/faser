package favicon

import (
	"golang.org/x/net/html"
	"os"
	"strings"
	"testing"
)

const testIndexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
	<title>Test</title>
	<link rel="icon" type="image/png" sizes="32x32" href="static/img/favicon-32x32.png" />
    <link rel="icon" type="image/png" sizes="16x16" href="static/img/favicon-16x16.png" />
</head>
<body>
	<p>Test</p>
</body>`

func TestCleanUpFiles(t *testing.T) {
	if err := os.MkdirAll("files/test", 0644); err != nil {
		t.Fatal(err)
	}

	cleanUpFiles("test")

	if _, err := os.Stat("files/test"); !os.IsNotExist(err) {
		t.Fatal("Directory must not exist anymore")
	}
}

func TestFindLinkIconNodes(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(testIndexHTML))

	if err != nil {
		t.Fatal(err)
	}

	head := findHeadNode(doc)

	if head == nil {
		t.Fatal("Head node must have been found")
	}

	links := findLinkIconNodes(head)

	if len(links) != 2 {
		t.Fatalf("Two link nodes must have been returned, but was: %v", len(links))
	}

	if getHref(&links[0]) != "static/img/favicon-32x32.png" ||
		getHref(&links[1]) != "static/img/favicon-16x16.png" {
		t.Fatalf("Returned link nodes not as expected: %v", links)
	}
}

func TestGetFaviconURL(t *testing.T) {
	in := []string{
		"favicon.ico",
		"favicon.ico?query=params&and=more",
		"some/path/favicon.ico",
		"http://hostname.com/some/path/favicon.ico",
		"https://hostname.com/some/path/favicon.ico",
	}
	out := []string{
		"http://hostname.com/favicon.ico",
		"http://hostname.com/favicon.ico?query=params&and=more",
		"http://hostname.com/some/path/favicon.ico",
		"http://hostname.com/some/path/favicon.ico",
		"https://hostname.com/some/path/favicon.ico",
	}

	for i, url := range in {
		if o := getFaviconURL(url, "hostname.com"); o != out[i] {
			t.Fatalf("Expected '%v' for '%v' but was '%v'", out[i], url, o)
		}
	}
}

func TestGetFaviconFilename(t *testing.T) {
	in := []string{
		"favicon.ico",
		"favicon.ico?query=params&and=more",
		"some/path/favicon.ico",
		"hostname.com/some/path/favicon.ico",
		"http://hostname.com/some/path/favicon.ico",
		"https://hostname.com/some/path/favicon.ico",
	}
	out := []string{
		"favicon.ico",
		"favicon.ico",
		"favicon.ico",
		"favicon.ico",
		"favicon.ico",
		"favicon.ico",
	}

	for i, url := range in {
		if o := getFaviconFilename(url); o != out[i] {
			t.Fatalf("Expected '%v' for '%v' but was '%v'", out[i], url, o)
		}
	}
}
