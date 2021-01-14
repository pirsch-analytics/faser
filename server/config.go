package server

import (
	"os"
	"strconv"
	"strings"
)

var cfg *config

type config struct {
	LogLevel        string
	CacheMaxAge     int
	CacheMaxEntries int
	DefaultFavicon  string
	Cors            corsConfig
	Server          serverConfig
	DB              postgresConfig
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
	MigrationDir       string
	SSLMode            string
	SSLCert            string
	SSLKey             string
	SSLRootCert        string
}

// LoadConfig loadsd the configuration from environment variables.
func LoadConfig() {
	cacheMaxAge, _ := strconv.Atoi(os.Getenv("FASER_CACHE_MAX_AGE"))
	cacheMaxEntries, _ := strconv.Atoi(os.Getenv("FASER_CACHE_MAX_ENTRIES"))
	defaultFavicon := os.Getenv("FASER_DEFAULT_FAVICON")
	writeTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_WRITE_TIMEOUT"))
	readTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_READ_TIMEOUT"))
	dbMaxOpenConnections, _ := strconv.Atoi(os.Getenv("FASER_DB_MAX_OPEN_CONNECTIONS"))

	if defaultFavicon == "" {
		defaultFavicon = "default.svg"
	}

	if cacheMaxAge <= 0 {
		cacheMaxAge = 3600 * 24 * 7 // one week in seconds
	}

	if cacheMaxEntries <= 0 {
		cacheMaxEntries = 10_000
	}

	if writeTimeout <= 0 {
		writeTimeout = 5
	}

	if readTimeout <= 0 {
		readTimeout = 5
	}

	cfg = &config{
		LogLevel:       strings.ToLower(os.Getenv("FASER_LOG_LEVEL")),
		CacheMaxAge:    cacheMaxAge,
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
			MigrationDir:       os.Getenv("FASER_DB_MIGRATION_DIR"),
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
