package main

import (
	"log"
	"net/http"
)

// watch:
//
//	find -name '*go' -or -name '*.html' | entr -crs "go run cmd/make/make.go build; go run cmd/make/make.go serve"
//
// build:
//
//	go build -gcflags "-m=1" -ldflags="-s" -o serve cmd/serve/serve.go
const (
	port    string = "8080"
	rootDir string = "."
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(rootDir))) // Serve files from project directory and subdirectories
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...", port, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
