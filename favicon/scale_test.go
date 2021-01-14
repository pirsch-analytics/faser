package favicon

import (
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
	files := []string{
		"test.gif",
		"test.ico",
		"test.jpg",
		"test.png",
	}

	for _, f := range files {
		scaledPath := filepath.Join(filesDir, "hostname", getFilenameForSize(f, 16))

		if _, err := os.Stat(scaledPath); !os.IsNotExist(err) {
			if err := os.Remove(scaledPath); err != nil {
				t.Fatal(err)
			}
		}

		if err := scale("hostname", f, 16); err != nil {
			t.Fatalf("File must have been rescaled, but was: %v", err)
		}

		if _, err := os.Stat(scaledPath); err != nil {
			t.Fatalf("Scaled image must exist, but was: %v", err)
		}
	}
}
