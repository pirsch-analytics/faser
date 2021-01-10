FROM golang AS build
RUN apt-get update && \
    apt-get upgrade -y
WORKDIR /go/src/faser
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/faser/main.go && \
	mkdir /app && \
	mv /go/src/faser/main /app/server && \
	mv /go/src/faser/schema /app/schema && \
	mv /go/src/faser/default.svg /app/default.svg

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 8080
VOLUME ["/app/files", "/app/default.svg"]
ENTRYPOINT ["/app/server"]
