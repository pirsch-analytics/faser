package favicon

import (
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logbuch.Debug(vars["url"])
	w.WriteHeader(http.StatusNotFound)
}
