#!/bin/bash

export FASER_LOG_LEVEL=debug
export FASER_CORS_ORIGINS=*
export FASER_SERVER_HOST=localhost:8080
go run main.go
