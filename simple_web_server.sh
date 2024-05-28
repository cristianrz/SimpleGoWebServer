#!/bin/sh

set -eu

_http_not_found() {
	printf "HTTP/1.1 404 Not Found\r\n"
	printf "Content-Type: text/html\r\n"
	printf "\r\n"
	printf "<html><body><p>404 Not Found</p></body></html>"
	exit 1
}

_handle_request() {
	read -r _method _request_path _version

	# decode
	_request_path="$(busybox httpd -d "$_request_path")"
	_request_path="$(echo "$_request_path" | sed 's/..\///g')"
	_request_path="${_request_path%%/}"

	while read -r _line; do
		_line="$(echo "$_line" | sed 's/\r//g')"
		[ "$_line" = "" ] && break
	done

	# Only handle GET requests
	if [ "$_method" != "GET" ]; then
		_http_not_found
	fi

	_filesystem_path="$DIR${_request_path%%/}"

	if [ "$_request_path" = "/" ]; then
		_request_path=""

		if [ -f "$_filesystem_path/index.html" ]; then
			_filesystem_path="${_filesystem_path}/index.html"
		fi
	fi

	# Check if it is a dir
	if [ -d "$_filesystem_path" ]; then
		# List the dir contents
		printf "HTTP/1.1 200 OK\r\n"
		printf "Content-Type: text/html\r\n"
		printf "\r\n"
		cat <<EOF
<!DOCTYPE HTML>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Directory listing for $_request_path</title>
<style>
body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
    background-color: #f6f8fa;
    margin: 0;
    padding: 0;
    color: #24292e;
}

h1 {
    font-size: 24px;
    color: #24292e;
    margin: 16px 0;
    padding: 16px;
    border-bottom: 1px solid #e1e4e8;
    background-color: #fff;
}

hr {
    border: 0;
    border-top: 1px solid #e1e4e8;
    margin: 0;
}

ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
}

li {
    display: flex;
    padding: 8px 16px;
    border-bottom: 1px solid #e1e4e8;
    background-color: #fff;
}

li:hover {
    background-color: #f6f8fa;
}

a {
    text-decoration: none;
    color: #0366d6;
    flex-grow: 1;
    /*font-weight: 600;*/
}

a:hover {
    text-decoration: underline;
}

.container {
    width: 80%;
    margin: 40px auto;
    background-color: #fff;
    border: 1px solid #e1e4e8;
    border-radius: 6px;
    box-shadow: 0 1px 3px rgba(27, 31, 35, 0.12);
}

</style>
</head>
<body>
<div class="container">
<h1>Directory listing for ${_request_path:-/}</h1>
<ul>
EOF
		printf '<li><a href="%s/..">⬆️</a></li>\n' "$(basename "$_request_path")"
		#shellcheck disable=SC2012
		ls -1 "$_filesystem_path" | awk -v path="$_request_path" '{print "<li><a href=\"" path "/" $0 "\">" $0 "</a></li>"}'
		cat <<EOF
  </ul>
</div>
</body>
</html>
EOF
	# Check if the file exists
	elif [ ! -e "$_filesystem_path" ]; then
		printf "HTTP/1.1 404 Not Found\r\n"
		printf "Content-Type: text/html\r\n"
		printf "\r\n"
		printf "<html><body><h1>404 Not Found</h1></body></html>"
		exit
	elif [ -x "$_filesystem_path" ]; then
		# Run the file as CGI script
		printf "HTTP/1.1 200 OK\r\n"
		printf "Content-Type: text/html\r\n"
		printf "\r\n"
		"$_filesystem_path"
	else
		# Serve the file as static HTML
		printf "HTTP/1.1 200 OK\r\n"
		printf "Content-Type: text/html\r\n"
		printf "\r\n"
		cat "$_filesystem_path"
	fi
}

_here="$(
	cd "$(dirname "$0")"
	pwd
)"

_this="$_here/$(basename "$0")"

if [ "${_script:-}" != "yes" ]; then
	export _script="yes"
	exec nc -l "$PORT" -k -e "$_this"
fi

_handle_request
