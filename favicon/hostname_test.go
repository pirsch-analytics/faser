package favicon

import "testing"

func TestGetHostname(t *testing.T) {
	in := []string{
		"example.com",
		"http://example.com",
		"https://example.com",
		"https://www.example.com",
		"example",
		"key words",
		"",
	}
	out := []string{
		"example.com",
		"example.com",
		"example.com",
		"www.example.com",
		"example",
		"",
		"",
	}

	for i, url := range in {
		if o := getHostname(url); o != out[i] {
			t.Fatalf("Expected '%v' for '%v' but was '%v'", out[i], url, o)
		}
	}
}
