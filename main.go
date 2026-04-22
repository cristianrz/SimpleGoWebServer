package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	port = os.Getenv("PORT")
	addr = os.Getenv("ADDR")
	dir  = os.Getenv("DIR")
)

const directoryListingStyle = `
<style>
body {
    font-family: 'Arial', sans-serif;
    background-color: #f9f9f9;
    color: #333;
    margin: 0;
    padding: 20px;
}
.container {
    width: 80%;
    max-width: 800px;
    margin: 0 auto;
    background: #fff;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
h1 {
    color: #444;
    font-size: 24px;
    border-bottom: 2px solid #eee;
    padding-bottom: 10px;
    margin-bottom: 20px;
}
ul { list-style-type: none; padding: 0; }
li { margin: 10px 0; }
a { text-decoration: none; color: #1a73e8; font-weight: bold; }
a:hover { text-decoration: underline; }
</style>
`

func main() {
	if port == "" {
		port = "8080"
	}
	if addr == "" {
		addr = "127.0.0.1"
	}
	if dir == "" {
		dir = "."
	}

	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving base directory: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handleRequest)

	listenAddr := fmt.Sprintf("%s:%s", addr, port)
	fmt.Printf("Listening on http://%s (serving %s)\n", listenAddr, dir)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(dir, filepath.FromSlash(r.URL.Path))
	path, err := filepath.Abs(path)
	if err != nil || (!strings.HasPrefix(path, dir+string(filepath.Separator)) && path != dir) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		handleDirectory(w, r, path)
	} else {
		handleFile(w, r, path)
	}
}

func handleDirectory(w http.ResponseWriter, r *http.Request, path string) {
	indexPath := filepath.Join(path, "index.html")
	info, err := os.Stat(indexPath)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, r, indexPath)
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><head>%s<title>Directory listing for %s</title></head><body><div class='container'><h1>Directory listing for %s</h1><ul>",
		directoryListingStyle, r.URL.Path, r.URL.Path)
	fmt.Fprintf(w, `<li><a href="..">⬆️ ..</a></li>`)

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		link := filepath.Join(r.URL.Path, name)
		if file.IsDir() {
			link += "/"
			name = "📁 " + name
		} else {
			info, err := file.Info()
			if err == nil && info.Mode()&0111 != 0 {
				name = "⚙️ " + name
			} else {
				name = "📄 " + name
			}
		}
		fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, link, name)
	}
	fmt.Fprintf(w, "</ul></div></body></html>")
}

func handleFile(w http.ResponseWriter, r *http.Request, path string) {
	if isExecutable(path) {
		handleCGI(w, r, path)
	} else {
		http.ServeFile(w, r, path)
	}
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// handleCGI runs the file as a CGI script. Scripts must write a complete HTTP
// response including headers separated from the body by a blank line:
//
//	Content-Type: text/plain
//
//	Hello, world!
func handleCGI(w http.ResponseWriter, r *http.Request, path string) {
	cmd := exec.Command(path)
	cmd.Stdin = r.Body

	env := os.Environ()
	env = append(env,
		fmt.Sprintf("REQUEST_METHOD=%s", r.Method),
		fmt.Sprintf("QUERY_STRING=%s", r.URL.RawQuery),
		fmt.Sprintf("CONTENT_TYPE=%s", r.Header.Get("Content-Type")),
		fmt.Sprintf("CONTENT_LENGTH=%d", r.ContentLength),
		fmt.Sprintf("SCRIPT_NAME=%s", r.URL.Path),
		fmt.Sprintf("SERVER_NAME=%s", r.Host),
		fmt.Sprintf("SERVER_PORT=%s", port),
		fmt.Sprintf("REMOTE_ADDR=%s", r.RemoteAddr),
	)
	for k, v := range r.Header {
		env = append(env, fmt.Sprintf("HTTP_%s=%s",
			strings.ToUpper(strings.ReplaceAll(k, "-", "_")), v[0]))
	}
	cmd.Env = env

	out, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Split headers from body on the first blank line
	parts := strings.SplitN(string(out), "\r\n\r\n", 2)
	if len(parts) == 1 {
		parts = strings.SplitN(string(out), "\n\n", 2)
	}

	if len(parts) == 2 {
		for _, line := range strings.Split(parts[0], "\n") {
			line = strings.TrimRight(line, "\r")
			if line == "" {
				continue
			}
			kv := strings.SplitN(line, ": ", 2)
			if len(kv) == 2 {
				w.Header().Set(kv[0], kv[1])
			}
		}
		fmt.Fprint(w, parts[1])
	} else {
		fmt.Fprint(w, string(out))
	}
}
