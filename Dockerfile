FROM golang AS build
RUN apt-get update && \
    apt-get upgrade -y
WORKDIR /go/src/faser
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/faser/main.go && \
	mkdir -p /app/data && \
	mv /go/src/faser/main /app/server && \
	mv /go/src/faser/default.svg /app/data/favicon.svg

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
ENV FASER_CACHE_DIR=/app/data/files
ENV FASER_DEFAULT_FAVICON=/app/data/favicon.svg

EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["/app/server"]
