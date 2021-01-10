package favicon

import "testing"

func TestGetDomain(t *testing.T) {
	in := []string{
		"example.com",
		"http://example.com",
		"https://example.com",
		"https://www.example.com",
		"example",
		"",
	}
	out := []string{
		"example.com",
		"example.com",
		"example.com",
		"www.example.com",
		"example",
		"",
	}

	for i, url := range in {
		if o := getDomain(url); o != out[i] {
			t.Fatalf("Expected '%v' for '%v' but was '%v'", out[i], url, o)
		}
	}
}
