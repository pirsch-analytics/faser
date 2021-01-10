package server

import (
	"os"
	"strconv"
	"strings"
)

var cfg *config

type config struct {
	LogLevel string
	Cors     corsConfig
	Server   serverConfig
}

type corsConfig struct {
	LogLevel string
	Origins  string
}

type serverConfig struct {
	Host         string
	WriteTimeout int
	ReadTimeout  int
	TLS          bool
	TLSCert      string
	TLSKey       string
}

// LoadConfig loadsd the configuration from environment variables.
func LoadConfig() {
	writeTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_WRITE_TIMEOUT"))
	readTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_READ_TIMEOUT"))

	if writeTimeout <= 0 {
		writeTimeout = 5
	}

	if readTimeout <= 0 {
		readTimeout = 5
	}

	cfg = &config{
		LogLevel: strings.ToLower(os.Getenv("FASER_LOG_LEVEL")),
		Cors: corsConfig{
			LogLevel: strings.ToLower(os.Getenv("FASER_CORS_LOG_LEVEL")),
			Origins:  os.Getenv("FASER_CORS_ORIGINS"),
		},
		Server: serverConfig{
			Host:         os.Getenv("FASER_SERVER_HOST"),
			WriteTimeout: writeTimeout,
			ReadTimeout:  readTimeout,
			TLS:          strings.ToLower(os.Getenv("FASER_SERVER_TLS")) == "true",
			TLSCert:      os.Getenv("FASER_SERVER_TLS_CERT"),
			TLSKey:       os.Getenv("FASER_SERVER_TLS_KEY"),
		},
	}
}

// Get returns the configuration.
func Get() *config {
	return cfg
}
