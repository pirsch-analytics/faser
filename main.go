package main

import (
	"context"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"github.com/pirsch-analytics/faser/db"
	"github.com/pirsch-analytics/faser/favicon"
	"github.com/pirsch-analytics/faser/server"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func startServer(handler http.Handler, shutdown func()) {
	c := server.Get()
	logbuch.Info("Starting server...")
	logbuch.Info("Using HTTP read/write timeouts", logbuch.Fields{"write_timeout": c.Server.WriteTimeout, "read_timeout": c.Server.ReadTimeout})

	s := &http.Server{
		Handler:      handler,
		Addr:         c.Server.Host,
		WriteTimeout: time.Duration(c.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.Server.ReadTimeout) * time.Second,
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logbuch.Info("Shutting down server...")

		if shutdown != nil {
			shutdown()
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

		if err := s.Shutdown(ctx); err != nil {
			logbuch.Fatal("Error shutting down server gracefully", logbuch.Fields{"err": err})
		}
	}()

	if c.Server.TLS {
		logbuch.Info("TLS enabled")

		if err := s.ListenAndServeTLS(c.Server.TLSCert, c.Server.TLSKey); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	} else {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	}
}

func main() {
	server.LoadConfig()
	server.ConfigureLogging()
	db.Migrate()
	db.Connect()
	router := mux.NewRouter()
	router.HandleFunc("/{url}", favicon.ServeFavicon)
	cors := server.ConfigureCors(router)
	startServer(cors, func() {
		db.Disconnect()
	})
}
