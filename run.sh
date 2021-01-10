#!/bin/bash

export FASER_LOG_LEVEL=debug
export FASER_CORS_ORIGINS=*
export FASER_SERVER_HOST=localhost:8080
export FASER_DB_HOST=localhost
export FASER_DB_PORT=5432
export FASER_DB_SCHEMA=faser
export FASER_DB_USER=postgres
export FASER_DB_PASSWORD=postgres
export FASER_DB_SSL_MODE=disable
go run main.go
