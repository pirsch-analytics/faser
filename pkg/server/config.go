package server

import (
	"os"
	"strconv"
	"strings"
)

var cfg *App

type App struct {
	LogLevel string
	Cache    Cache
	Cors     CORS
	Server   Server
}

type Cache struct {
	Dir               string
	MaxAge            int
	MaxEntries        int
	DefaultFavicon    string
	DefaultFaviconDir string
}

type CORS struct {
	LogLevel string
	Origins  string
}

type Server struct {
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
	defaultFaviconDir := os.Getenv("FASER_DEFAULT_FAVICON_DIR")
	writeTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_WRITE_TIMEOUT"))
	readTimeout, _ := strconv.Atoi(os.Getenv("FASER_SERVER_READ_TIMEOUT"))

	if cacheDir == "" {
		cacheDir = "cache"
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

	if defaultFaviconDir == "" {
		defaultFaviconDir = "static"
	}

	if writeTimeout <= 0 {
		writeTimeout = 5
	}

	if readTimeout <= 0 {
		readTimeout = 5
	}

	cfg = &App{
		LogLevel: strings.ToLower(os.Getenv("FASER_LOG_LEVEL")),
		Cache: Cache{
			Dir:               cacheDir,
			MaxAge:            cacheMaxAge,
			MaxEntries:        cacheMaxEntries,
			DefaultFavicon:    defaultFavicon,
			DefaultFaviconDir: defaultFaviconDir,
		},
		Cors: CORS{
			LogLevel: strings.ToLower(os.Getenv("FASER_CORS_LOG_LEVEL")),
			Origins:  os.Getenv("FASER_CORS_ORIGINS"),
		},
		Server: Server{
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
func Config() *App {
	return cfg
}
