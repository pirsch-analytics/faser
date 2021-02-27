package favicon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeFaviconBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	ServeFavicon(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatal("Request must return a bad request status code")
	}

	body, _ := io.ReadAll(w.Body)

	if !strings.Contains(string(body), "provide a valid URL or hostname") ||
		strings.Contains(string(body), "provide a number greater or equal to 0") {
		t.Fatalf("Body must contain an error message, but was: %v", string(body))
	}
}

func TestServeFaviconBadRequestSize(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?url=example.com&size=asdf", nil)
	ServeFavicon(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatal("Request must return a bad request status code")
	}

	body, _ := io.ReadAll(w.Body)

	if !strings.Contains(string(body), "provide a valid URL or hostname") ||
		!strings.Contains(string(body), "provide a number greater or equal to 0") {
		t.Fatalf("Body must contain an error message, but was: %v", string(body))
	}
}
