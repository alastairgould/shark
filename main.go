package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

var packSizes = []int{250, 500}

type packRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	Packs map[int]int `json:"packs"`
}

func calculatePacks(quantity int) map[int]int {
	if quantity <= 0 || len(packSizes) == 0 {
		return map[int]int{}
	}

	ordered := append([]int(nil), packSizes...)
	sort.Sort(sort.Reverse(sort.IntSlice(ordered)))
	smallest := ordered[len(ordered)-1]

	total := ((quantity + smallest - 1) / smallest) * smallest

	packs := map[int]int{}
	for _, size := range ordered {
		if n := total / size; n > 0 {
			packs[size] = n
			total -= n * size
		}
	}
	return packs
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pack", func(w http.ResponseWriter, r *http.Request) {
		var req packRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp := packResponse{Packs: calculatePacks(req.Quantity)}

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
