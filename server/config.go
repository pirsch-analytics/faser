package server

import (
	"os"
	"strconv"
	"strings"
)

var cfg *config

type config struct {
	LogLevel       string
	Cache          int
	DefaultFavicon string
	Cors           corsConfig
	Server         serverConfig
	DB             postgresConfig
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

type postgresConfig struct {
	Host               string
	Port               string
	User               string
	Password           string
	Schema             string
	MaxOpenConnections int
	SSLMode            string
	SSLCert            string
	SSLKey             string
	SSLRootCert        string
}

// LoadConfig loadsd the configuration from environment variables.
func LoadConfig() {
	cache, _ := strconv.Atoi(os.Getenv("FASER_CACHE"))
	defaultFavicon := os.Getenv("FASER_DEFAULT_FAVICON")
	writeTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_WRITE_TIMEOUT"))
	readTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_READ_TIMEOUT"))
	dbMaxOpenConnections, _ := strconv.Atoi(os.Getenv("FASER_DB_MAX_OPEN_CONNECTIONS"))

	if defaultFavicon == "" {
		defaultFavicon = "default.svg"
	}

	if cache <= 0 {
		cache = 3600 * 24 * 7 // one week in seconds
	}

	if writeTimeout <= 0 {
		writeTimeout = 5
	}

	if readTimeout <= 0 {
		readTimeout = 5
	}

	cfg = &config{
		LogLevel:       strings.ToLower(os.Getenv("FASER_LOG_LEVEL")),
		Cache:          cache,
		DefaultFavicon: defaultFavicon,
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
		DB: postgresConfig{
			Host:               os.Getenv("FASER_DB_HOST"),
			Port:               os.Getenv("FASER_DB_PORT"),
			User:               os.Getenv("FASER_DB_USER"),
			Password:           os.Getenv("FASER_DB_PASSWORD"),
			Schema:             os.Getenv("FASER_DB_SCHEMA"),
			MaxOpenConnections: dbMaxOpenConnections,
			SSLMode:            os.Getenv("FASER_DB_SSL_MODE"),
			SSLCert:            os.Getenv("FASER_DB_SSL_CERT"),
			SSLKey:             os.Getenv("FASER_DB_SSL_KEY"),
			SSLRootCert:        os.Getenv("FASER_DB_SSL_ROOT_CERT"),
		},
	}
}

// Config returns the configuration.
func Config() *config {
	return cfg
}
