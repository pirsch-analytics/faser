# Faser

[![Go Report Card](https://goreportcard.com/badge/github.com/pirsch-analytics/faser)](https://goreportcard.com/report/github.com/pirsch-analytics/faser)

Faser (*fa*vicon *ser*ver) is a self-hosted, cached favicon server that returns an image for a domain or URL.

## Installation

Please see the docker-compose.yml example below for reference. You can set the following environment variables:

| Variable                   | Description                                                                 |
|----------------------------|-----------------------------------------------------------------------------|
| FASER_LOG_LEVEL            | debug, info, warn                                                           |
| FASER_CACHE_DIR            | Sets the favicon cache directory. `/app/cache` inside the container volume. |
| FASER_CACHE_MAX_AGE        | Sets the maximum time a favicon will be cached (in seconds).                |
| FASER_CACHE_MAX_ENTRIES    | Sets the maximum number of favicons stored (not including scaled images).   |
| FASER_DEFAULT_FAVICON      | Sets the default favicon path. `default.svg` inside the container volume.   |
| FASER_DEFAULT_FAVICON_DIR  | Sets the fallback favicon path. `/app/static` inside the container volume.  |
| FASER_CORS_LOG_LEVEL       | debug, info                                                                 |
| FASER_CORS_ORIGINS         | Sets the allowed origins (`*` by default).                                  |
| FASER_SERVER_HOST          | Sets the host ip:port (`:8080` by default).                                 |
| FASER_SERVER_WRITE_TIMEOUT | Sets the HTTP write timeout in seconds (5 by default).                      |
| FASER_SERVER_READ_TIMEOUT  | Sets the HTTP read timeout in seconds (5 by default).                       |
| FASER_SERVER_TLS           | true, false                                                                 |
| FASER_SERVER_TLS_CERT      | Path for the TLS certificate.                                               |
| FASER_SERVER_TLS_KEY       | Path for the TLS private key.                                               |

## Usage

Once you have the server up and running, you can make a request to it:

```
example.com/?url=github.com&size=64
```

This request will return the favicon for `github.com` scaled to a maximum of 64px width and height. The `size` parameter is optional. If you don't provide it, the maximum size will be returned. Make sure that you still set the maximum sizes on your end, as the favicon might not have the exact size you specified (if it's a svg for example).

In case the favicon cannot be downloaded, the default favicon is returned. You can also select a favicon if you know its filename by adding the `fallback` query parameter.

```
example.com/?url=github.com&size=64&fallback=my-favicon.svg
```

## Example Docker Compose

```
version: "3"

services:
  faser:
    image: pirsch/faser
    container_name: faser
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
