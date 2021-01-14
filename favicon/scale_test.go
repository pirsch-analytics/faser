package favicon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSelectFilenameForSize(t *testing.T) {
	if o, s := selectFilenameForSize("foo.png", 12); o != "foo-16.png" || s != 16 {
		t.Fatalf("Filename with minimum size must have been returned, but was: %v %v", o, s)
	}

	if o, s := selectFilenameForSize("foo.PNG", 1218); o != "foo-196.png" || s != 196 {
		t.Fatalf("Filename with maximum size must have been returned, but was: %v %v", o, s)
	}

	if o, s := selectFilenameForSize("foo.png", 64); o != "foo-64.png" || s != 64 {
		t.Fatalf("Filename with exact size must have been returned, but was: %v %v", o, s)
	}

	if o, s := selectFilenameForSize("foo.svg", 64); o != "foo.svg" || s != 0 {
		t.Fatalf("Original filename must have been returned, but was: %v %v", o, s)
	}
}

func TestScale(t *testing.T) {
	if err := os.MkdirAll("files/hostname", 0744); err != nil {
		t.Fatal(err)
	}

	files := []string{
		"test.gif",
		"test.ico",
		"test.jpg",
		"test.png",
	}

	for _, f := range files {
		originalPath := filepath.Join("test", f)
		targetPath := filepath.Join("files", "hostname", f)

		if _, err := os.Stat(targetPath); err != nil {
			copyFile(t, originalPath, targetPath)
		}

		scaledPath := filepath.Join("files", "hostname", getFilenameForSize(f, 16))

		if _, err := os.Stat(scaledPath); !os.IsNotExist(err) {
			if err := os.Remove(scaledPath); err != nil {
				t.Fatal(err)
			}
		}

		if filename := scale("hostname", f, 16); filename != getFilenameForSize(f, 16) {
			t.Fatalf("File must have been rescaled, but was: %v", filename)
		}

		if _, err := os.Stat(scaledPath); err != nil {
			t.Fatalf("Scaled image must exist, but was: %v", err)
		}
	}
}

func copyFile(t *testing.T, src, target string) {
	data, err := ioutil.ReadFile(src)

	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(target, data, 0744); err != nil {
		t.Fatal(err)
	}
}
