package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type packRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	Packs map[int]int `json:"packs"`
}

type problemDetail struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(problemDetail{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}); err != nil {
		log.Printf("encode problem: %v", err)
	}
}

func handler(p *packer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pack", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var req packRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeProblem(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Quantity < 1 {
			writeProblem(w, http.StatusBadRequest, "quantity must be at least 1")
			return
		}
		if req.Quantity > p.maxQuantity {
			writeProblem(w, http.StatusBadRequest, fmt.Sprintf("quantity must not exceed %d", p.maxQuantity))
			return
		}

		resp := packResponse{Packs: p.calculate(req.Quantity)}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("encode response: %v", err)
		}
	})
	return mux
}
