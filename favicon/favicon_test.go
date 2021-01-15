package favicon

import (
	"golang.org/x/net/html"
	"os"
	"strings"
	"testing"
)

const (
	testIndexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <base href="/" />
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="copyright" content="Marvin Blum" />
    <meta name="author" content="Marvin Blum" />
    <meta name="title" content="Marvin Blum" />
    <meta name="description" content="A full stack software engineer from Germany, open source and Linux enthusiast and co-founder of Emvi." />
    <meta name="msapplication-TileColor" content="#000000" />
    <meta name="theme-color" content="#000000" />
    <meta name="twitter:card" content="profile" />
    <meta name="twitter:site" content="@m5blum" />
    <meta name="twitter:title" content="Marvin Blum" />
    <meta name="twitter:description" content="A full stack software engineer from Germany, open source and Linux enthusiast and co-founder of Emvi." />
    <meta name="twitter:image" content="https://marvinblum.de/avatar.png" />
    <meta property="og:url" content="https://marvinblum.de/" />
    <meta property="og:title" content="Marvin Blum" />
    <meta property="og:description" content="A full stack software engineer from Germany, open source and Linux enthusiast and co-founder of Emvi." />
    <meta property="og:image" content="https://marvinblum.de/avatar.png" />
    <title>marvin blum</title>
    <link rel="apple-touch-icon" sizes="180x180" href="../static/favicon/apple-touch-icon.png" />
    <link rel="icon" type="image/png" sizes="32x32" href="../static/favicon/favicon-32x32.png" />
    <link rel="icon" type="image/png" sizes="16x16" href="../static/favicon/favicon-16x16.png" />
    <link rel="manifest" href="../static/favicon/site.webmanifest" />
    <link rel="mask-icon" href="../static/favicon/safari-pinned-tab.svg" color="#5bbad5" />
    <meta name="msapplication-TileColor" content="#da532c" />
    <meta name="theme-color" content="#ffffff" />
    <link rel="stylesheet" type="text/css" href="../static/normalize.css" />
    <link rel="stylesheet" type="text/css" href="../static/concrete.css" />
    <link rel="stylesheet" type="text/css" href="../static/style.css" />
</head>
<body>
	...
</body>
</html>`
)

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

	if len(links) != 3 {
		t.Fatalf("Three link nodes must have been returned, but was: %v", len(links))
	}

	if getHref(&links[0]) != "../static/favicon/apple-touch-icon.png" ||
		getHref(&links[1]) != "../static/favicon/favicon-32x32.png" ||
		getHref(&links[2]) != "../static/favicon/favicon-16x16.png" {
		t.Fatalf("Returned link nodes not as expected: %v", links)
	}

	largest := findLargestIconNode(links)

	if largest == nil || getHref(largest) != "../static/favicon/apple-touch-icon.png" {
		t.Fatalf("Largest icon must have been returned, but was: %v", getHref(largest))
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
