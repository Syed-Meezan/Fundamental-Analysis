// Command sampleproject runs a fundamental-analysis web app for Indian
// listed companies: it scrapes screener.in for financial statements,
// computes growth-trend metrics in Go, and serves the results through a
// small HTTP API plus a static frontend.
package main

import (
	"log"
	"net/http"

	"sampleproject/internal/httpapi"
	"sampleproject/web"
)

const addr = ":8080"

func main() {
	mux := http.NewServeMux()

	api := httpapi.NewServer()
	api.Routes(mux)

	mux.Handle("/static/", http.FileServerFS(web.FS))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.FS, "static/index.html")
	})

	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
