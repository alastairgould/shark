package main

import (
	_ "embed"
	"io"
	"log"
	"net/http"
)

//go:embed index.html
var indexHTML string

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := io.WriteString(w, indexHTML); err != nil {
		log.Printf("write page: %v", err)
	}
}
