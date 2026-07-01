package main

import "slices"

type packer struct {
	maxQuantity int
	packsFor    []int
	lastSize    []int
}

func precomputePackingTable(sizes []int, maxQuantity int) *packer {
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

func (p *packer) calculate(quantity int) map[int]int {
	total := quantity
	for total < len(p.packsFor) && p.packsFor[total] == -1 {
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
