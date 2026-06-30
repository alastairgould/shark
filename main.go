package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var defaultPackSizes = []int{250, 500, 1000, 2000, 5000}

type packRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	Packs map[int]int `json:"packs"`
}

func calculatePacks(quantity int, sizes []int) map[int]int {
	if quantity <= 0 || len(sizes) == 0 {
		return map[int]int{}
	}

	ordered := append([]int(nil), sizes...)
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

func handler(packSizes []int) http.Handler {
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

// packSizesFromEnv reads PACK_SIZES as a comma-separated list, falling back to
// defaultPackSizes when unset or invalid.
func packSizesFromEnv() []int {
	raw := os.Getenv("PACK_SIZES")
	if raw == "" {
		return defaultPackSizes
	}

	var sizes []int
	for field := range strings.SplitSeq(raw, ",") {
		size, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil || size <= 0 {
			log.Printf("ignoring invalid PACK_SIZES value %q, using defaults", raw)
			return defaultPackSizes
		}
		sizes = append(sizes, size)
	}
	return sizes
}

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler(packSizesFromEnv())); err != nil {
		log.Fatal(err)
	}
}
