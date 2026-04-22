# SimpleGoWebServer

A web server where any executable file is a CGI endpoint.

Drop a shell script in a directory, make it executable, point
SimpleGoWebServer at it — it's now an HTTP API. No config, no framework,
no routes to define.

```sh
DIR=/var/www PORT=8080 ADDR=0.0.0.0 simplegowebserver
```

## How it works

- **Executable file** → runs as a CGI script, output served as HTTP response
- **Regular file** → served as-is
- **Directory** → serves `index.html` if present, otherwise a directory listing
- **Hidden files** → never served

## CGI scripts

Any executable in the served directory becomes an endpoint. The script receives
standard CGI environment variables and must write a response with headers:

```sh
#!/bin/sh
# api/hello — available at /api/hello
printf 'Content-Type: text/plain\n\nHello, %s!\n' "$QUERY_STRING"
```

Available environment variables:

| Variable | Description |
|---|---|
| `REQUEST_METHOD` | GET, POST, etc. |
| `QUERY_STRING` | URL query string |
| `CONTENT_TYPE` | Request content type |
| `CONTENT_LENGTH` | Request body length |
| `SCRIPT_NAME` | Request path |
| `SERVER_NAME` | Host header |
| `SERVER_PORT` | Listening port |
| `REMOTE_ADDR` | Client address |
| `HTTP_*` | All request headers |

## Install

```sh
go install github.com/cristianrz/SimpleGoWebServer@latest
```

Or build from source:

```sh
git clone https://github.com/cristianrz/SimpleGoWebServer.git
cd SimpleGoWebServer
go build -o simplegowebserver .
```

## Usage

```sh
# Serve current directory on localhost:8080
simplegowebserver

# Serve a specific directory on all interfaces
DIR=/var/www PORT=80 ADDR=0.0.0.0 simplegowebserver
```

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port to listen on |
| `ADDR` | `127.0.0.1` | Address to bind |
| `DIR` | `.` | Directory to serve |

## Docker

```sh
docker compose up
```

## License

MIT
