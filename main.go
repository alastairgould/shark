package main

import (
	"log"
	"net/http"
)

// handler returns the HTTP handler for the application.
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pack", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func main() {
	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler()); err != nil {
		log.Fatal(err)
	}
}
