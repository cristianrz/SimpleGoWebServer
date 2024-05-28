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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r)
	})

	listenAddr := fmt.Sprintf("%s:%s", addr, port)
	fmt.Printf("Listening on %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(dir, r.URL.Path)

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
	fmt.Fprintf(w, "<html><body><h1>Directory listing for %s</h1><ul>", r.URL.Path)
	fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, "..", "⬆️")
	for _, file := range files {
		name := file.Name()
		link := filepath.Join(r.URL.Path, name)
		if file.IsDir() {
			link += "/"
		}
		fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, link, name)
	}
	fmt.Fprintf(w, "</ul></body></html>")
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
	mode := info.Mode()
	return mode&0111 != 0
}

func handleCGI(w http.ResponseWriter, r *http.Request, path string) {
	cmd := exec.Command(path)
	cmd.Stdin = r.Body
	cmd.Stdout = w
	cmd.Stderr = w

	env := os.Environ()
	env = append(env, fmt.Sprintf("REQUEST_METHOD=%s", r.Method))
	env = append(env, fmt.Sprintf("QUERY_STRING=%s", r.URL.RawQuery))
	env = append(env, fmt.Sprintf("CONTENT_TYPE=%s", r.Header.Get("Content-Type")))
	env = append(env, fmt.Sprintf("CONTENT_LENGTH=%d", r.ContentLength))

	for k, v := range r.Header {
		env = append(env, fmt.Sprintf("HTTP_%s=%s", strings.ToUpper(strings.ReplaceAll(k, "-", "_")), v[0]))
	}

	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

