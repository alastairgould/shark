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

func handler(p *packer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pack", func(w http.ResponseWriter, r *http.Request) {
		var req packRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Quantity < 1 {
			http.Error(w, "quantity must be at least 1", http.StatusBadRequest)
			return
		}
		if req.Quantity > p.maxQuantity {
			http.Error(w, fmt.Sprintf("quantity must not exceed %d", p.maxQuantity), http.StatusBadRequest)
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
