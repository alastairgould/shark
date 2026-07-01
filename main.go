package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var defaultPackSizes = []int{250, 500, 1000, 2000, 5000}

const defaultMaxQuantity = 1_000_000

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	p := precomputePackingTable(packSizesFromEnv(), maxQuantityFromEnv())

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler(p),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on %s", addr)
	return srv.ListenAndServe()
}

func handler(p *packer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleIndex)
	mux.HandleFunc("POST /pack", handlePack(p))
	return mux
}

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
