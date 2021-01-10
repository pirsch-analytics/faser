package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strings"
)

// ConfigureCors configures CORS log level and origins.
// The configuration is not restricted by default.
func ConfigureCors(router *mux.Router) http.Handler {
	origins := strings.Split(cfg.Cors.Origins, ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            cfg.Cors.LogLevel == "debug",
	})
	return c.Handler(router)
}
