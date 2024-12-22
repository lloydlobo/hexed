package main

import (
	"log"
	"log/slog"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	slog.Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
