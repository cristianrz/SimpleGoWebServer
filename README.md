# SimpleGoWebServer

SimpleGoWebServer is a lightweight web server written in Go that serves files
and directories from a specified base directory. It handles directories by
serving an index.html file if present, providing a directory listing if not,
and executes executable files as CGI scripts.

## Features

- Serve Static Files: Serves files as-is from the specified base directory.
- Directory Listing: Provides a directory listing if no index.html is present.
- CGI Script Execution: Executes executable files as CGI scripts.

## Prerequisites

- Go (1.16+)

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/SimpleGoWebServer.git
cd SimpleGoWebServer
```

Build the server:

```bash
go build -o simplegowebserver main.go
```

## Usage

Set the necessary environment variables:

- `PORT`: The port on which the server will listen (default: 8080).
- `ADDR`: The address on which the server will listen (default: 127.0.0.1).
- `DIR`: The base directory from which files will be served (default: current
  directory).

Example:

```bash
PORT=8080 ADDR=0.0.0.0 DIR=/path/to/directory ./simplegowebserver
```

Access the server in your browser:

```
http://<ADDR>:<PORT>
```

## Why

I wanted a simple oneliner that has CGI on executable and directory listing for
quickly spinning up any type of web server including APIs.

## License

This project is licensed under the MIT License. See the LICENSE file
for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for
any improvements or bug fixes.

