package favicon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeFavicon(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	ServeFavicon(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Request must return status code 200, but was %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Body)

	if !strings.Contains(string(body), "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" fill=\"none\"><path d=\"M21 12a9 9 0 0 1-9 9m9-9a9 9 0 0 0-9-9m9 9H3m9 9a9 9 0 0 1-9-9m9 9c2.501-2.465 3.923-5.663 4-9-.077-3.337-1.499-6.535-4-9m0 18c-2.501-2.465-3.923-5.663-4-9 .077-3.337 1.499-6.535 4-9m-9 9a9 9 0 0 1 9-9\" stroke=\"#000\" stroke-width=\"2\"/></svg>\n") {
		t.Fatalf("Body must contain default favicon, but was: %v", string(body))
	}
}
