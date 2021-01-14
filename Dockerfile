FROM golang AS build
RUN apt-get update && \
    apt-get upgrade -y
WORKDIR /go/src/faser
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/faser/main.go && \
	mkdir /app && \
	mkdir /app/default && \
	mv /go/src/faser/main /app/server && \
	mv /go/src/faser/default.svg /app/default/favicon.svg

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates imagemagick && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

ENV FASER_LOG_LEVEL=info
ENV FASER_CORS_ORIGINS=*
ENV FASER_SERVER_HOST=:8080
ENV FASER_DEFAULT_FAVICON=/app/default/favicon.svg

EXPOSE 8080
VOLUME ["/app/files", "/app/default"]
ENTRYPOINT ["/app/server"]
