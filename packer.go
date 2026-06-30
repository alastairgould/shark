package main

import "slices"

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
