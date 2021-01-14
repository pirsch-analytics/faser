package server

import (
	"os"
	"strconv"
	"strings"
)

var cfg *config

type config struct {
	LogLevel string
	Cache    cacheConfig
	Cors     corsConfig
	Server   serverConfig
}

type cacheConfig struct {
	Dir            string
	MaxAge         int
	MaxEntries     int
	DefaultFavicon string
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
	cacheDir := os.Getenv("FASER_CACHE_DIR")
	cacheMaxAge, _ := strconv.Atoi(os.Getenv("FASER_CACHE_MAX_AGE"))
	cacheMaxEntries, _ := strconv.Atoi(os.Getenv("FASER_CACHE_MAX_ENTRIES"))
	defaultFavicon := os.Getenv("FASER_DEFAULT_FAVICON")
	writeTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_WRITE_TIMEOUT"))
	readTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_READ_TIMEOUT"))

	if cacheDir == "" {
		cacheDir = "files"
	}

	if cacheMaxAge <= 0 {
		cacheMaxAge = 3600 * 24 * 7 // one week in seconds
	}

	if cacheMaxEntries <= 0 {
		cacheMaxEntries = 10_000
	}

	if defaultFavicon == "" {
		defaultFavicon = "default.svg"
	}

	if writeTimeout <= 0 {
		writeTimeout = 5
	}

	if readTimeout <= 0 {
		readTimeout = 5
	}

	cfg = &config{
		LogLevel: strings.ToLower(os.Getenv("FASER_LOG_LEVEL")),
		Cache: cacheConfig{
			Dir:            cacheDir,
			MaxAge:         cacheMaxAge,
			MaxEntries:     cacheMaxEntries,
			DefaultFavicon: defaultFavicon,
		},
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

// Config returns the configuration.
func Config() *config {
	return cfg
}
