package server

import (
	"context"
	"github.com/emvi/logbuch"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Start starts the server for given router and calls the function to gracefully shut down the server.
func Start(handler http.Handler, shutdown func()) {
	logbuch.Info("Starting server...")
	logbuch.Info("Using HTTP read/write timeouts", logbuch.Fields{
		"write_timeout": cfg.Server.WriteTimeout,
		"read_timeout":  cfg.Server.ReadTimeout,
	})

	s := &http.Server{
		Handler:      handler,
		Addr:         cfg.Server.Host,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
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

	if cfg.Server.TLS {
		logbuch.Info("TLS enabled")

		if err := s.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	} else {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	}
}
