package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

var defaultPackSizes = []int{250, 500, 1000, 2000, 5000}

const defaultMaxQuantity = 1_000_000

type packRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	Packs map[int]int `json:"packs"`
}

// packer holds a pack-size configuration and a precomputed table mapping every
// item total (up to maxQuantity + the largest pack) to the fewest packs that
// make it. The table depends only on the pack sizes, so it is built once and
// reused across requests; it is immutable after construction and safe for
// concurrent use.
type packer struct {
	maxQuantity int
	packsFor    []int
	lastSize    []int
}

// newPacker precomputes the packing table for the given sizes, supporting
// orders up to maxQuantity.
func newPacker(sizes []int, maxQuantity int) *packer {
	upper := maxQuantity + slices.Max(sizes)

	packsFor := make([]int, upper+1)
	lastSize := make([]int, upper+1)
	for itemTotal := 1; itemTotal <= upper; itemTotal++ {
		packsFor[itemTotal] = -1
	}

	for itemTotal := 1; itemTotal <= upper; itemTotal++ {
		for _, packSize := range sizes {
			remainder := itemTotal - packSize
			if packSize > itemTotal || packsFor[remainder] == -1 {
				continue
			}
			if packsFor[itemTotal] == -1 || packsFor[remainder]+1 < packsFor[itemTotal] {
				packsFor[itemTotal] = packsFor[remainder] + 1
				lastSize[itemTotal] = packSize
			}
		}
	}

	return &packer{maxQuantity: maxQuantity, packsFor: packsFor, lastSize: lastSize}
}

// calculate returns the packs to ship for an order. It reconstructs the answer
// from the precomputed table, so it only does the per-request work of finding
// the fewest-items total and walking it back into packs. The caller must pass a
// quantity in the range [1, maxQuantity]; the handler validates this.
func (p *packer) calculate(quantity int) map[int]int {
	total := quantity
	for p.packsFor[total] == -1 {
		total++
	}

	packs := map[int]int{}
	for total > 0 {
		packSize := p.lastSize[total]
		packs[packSize]++
		total -= packSize
	}
	return packs
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

// maxQuantityFromEnv reads MAX_QUANTITY, falling back to defaultMaxQuantity when
// unset or invalid.
func maxQuantityFromEnv() int {
	raw := os.Getenv("MAX_QUANTITY")
	if raw == "" {
		return defaultMaxQuantity
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		log.Printf("ignoring invalid MAX_QUANTITY value %q, using default", raw)
		return defaultMaxQuantity
	}
	return n
}

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	p := newPacker(packSizesFromEnv(), maxQuantityFromEnv())

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler(p)); err != nil {
		log.Fatal(err)
	}
}
