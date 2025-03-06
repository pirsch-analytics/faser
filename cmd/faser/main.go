package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/pirsch-analytics/faser/pkg/favicon"
	"github.com/pirsch-analytics/faser/pkg/server"
	"net/http"
	"strings"
)

func main() {
	server.LoadConfig()
	server.ConfigureLogging()
	favicon.Init()
	router := chi.NewRouter()
	cfg := server.Config()
	origins := strings.Split(cfg.Cors.Origins, ",")
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{http.MethodGet},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            cfg.Cors.LogLevel == "debug",
	}))
	router.MethodFunc(http.MethodGet, "/", favicon.ServeFavicon)
	server.Start(router, nil)
}
