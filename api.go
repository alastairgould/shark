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

func validationError(w http.ResponseWriter, detail string, code int) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(problemDetail{
		Type:   "about:blank",
		Title:  http.StatusText(code),
		Status: code,
		Detail: detail,
	}); err != nil {
		log.Printf("encode problem: %v", err)
	}
}

func handlePack(p *packer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var req packRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			validationError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Quantity < 1 {
			validationError(w, "quantity must be at least 1", http.StatusBadRequest)
			return
		}
		if req.Quantity > p.maxQuantity {
			validationError(w, fmt.Sprintf("quantity must not exceed %d", p.maxQuantity), http.StatusBadRequest)
			return
		}

		resp := packResponse{Packs: p.calculate(req.Quantity)}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}
