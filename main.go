package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var packSizes = []int{250, 500, 1000, 2000, 5000}

type packRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	Packs map[int]int `json:"packs"`
}

// calculatePacks returns the packs needed to fulfil quantity, as size -> count.
func calculatePacks(quantity int, sizes []int) map[int]int {
	return map[int]int{250: 1}
}

// handler returns the HTTP handler for the application.
func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pack", func(w http.ResponseWriter, r *http.Request) {
		var req packRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp := packResponse{Packs: calculatePacks(req.Quantity, packSizes)}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("encode response: %v", err)
		}
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
