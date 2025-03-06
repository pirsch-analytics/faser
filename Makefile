.PHONY: dev test deps release

dev:
	FASER_LOG_LEVEL=debug FASER_CORS_ORIGINS=* FASER_SERVER_HOST=localhost:8080 go run cmd/faser/main.go

test:
	go test -cover -race github.com/pirsch-analytics/faser/pkg/...

deps:
	go get -u -t ./...
	go mod tidy
	go mod vendor

docker:
	docker build -t pirsch/faser:$(VERSION) -f build/Dockerfile .
	docker push pirsch/faser:$(VERSION)
