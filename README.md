# URL Shortener

A simple URL shortener service built using Go and Redis.

## Features

* Rate Limit by IP

## Getting Started

### Local

```
Usage: url-shortener [options]
```

### Options

```
--address, -a           IP address to listen (default: "127.0.0.1")
--port, -p              Port number to listen (default: "8080")
--redis-addr, -r        Redis address (default: "localhost:6379")
--redis-password, -rp   Redis password (default: "")
```

### Docker

```
docker compose up
```
