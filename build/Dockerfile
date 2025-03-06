FROM golang AS build
RUN apt-get update && \
    apt-get upgrade -y
WORKDIR /go/src/faser
COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/faser/cmd/faser/main.go && \
	mkdir -p /app/data && \
	mv /go/src/faser/main /app/server && \
	mv /go/src/faser/static /app/static

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates imagemagick && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
WORKDIR /app

ENV FASER_LOG_LEVEL=info
ENV FASER_CORS_ORIGINS=*
ENV FASER_SERVER_HOST=:8080
ENV FASER_CACHE_DIR=/app/cache
ENV FASER_DEFAULT_FAVICON_DIR=/app/static

EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["/app/server"]
