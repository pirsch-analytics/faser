package main

import (
	"github.com/gorilla/mux"
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/favicon"
	"github.com/pirsch-analytics/faser/server"
)

func main() {
	server.LoadConfig()
	server.ConfigureLogging()
	db.Migrate()
	db.Connect()
	favicon.Init()
	router := mux.NewRouter()
	router.HandleFunc("/", favicon.ServeFavicon)
	cors := server.ConfigureCors(router)
	server.Start(cors, func() {
		db.Disconnect()
	})
}
