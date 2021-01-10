package favicon

import "testing"

func TestGetFilenameForSize(t *testing.T) {
	if o := getFilenameForSize("foo.png", 12); o != "foo-16.png" {
		t.Fatalf("Filename with minimum size must have been returned, but was: %v", o)
	}

	if o := getFilenameForSize("foo.png", 1218); o != "foo-196.png" {
		t.Fatalf("Filename with maximum size must have been returned, but was: %v", o)
	}

	if o := getFilenameForSize("foo.png", 64); o != "foo-64.png" {
		t.Fatalf("Filename with exact size must have been returned, but was: %v", o)
	}
}
