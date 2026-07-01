package main

import (
	"reflect"
	"testing"
)

func TestPackSizesFromEnv(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want []int
	}{
		{"unset falls back to defaults", "", defaultPackSizes},
		{"parses a comma-separated list", "250,500,1000", []int{250, 500, 1000}},
		{"trims surrounding whitespace", " 250 , 500 ", []int{250, 500}},
		{"non-numeric falls back to defaults", "abc", defaultPackSizes},
		{"non-positive falls back to defaults", "250,-5", defaultPackSizes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PACK_SIZES", tt.env)

			got := packSizesFromEnv()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("packSizesFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxQuantityFromEnv(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want int
	}{
		{"unset falls back to default", "", defaultMaxQuantity},
		{"parses a valid number", "5000", 5000},
		{"non-numeric falls back to default", "abc", defaultMaxQuantity},
		{"zero falls back to default", "0", defaultMaxQuantity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MAX_QUANTITY", tt.env)

			got := maxQuantityFromEnv()
			if got != tt.want {
				t.Errorf("maxQuantityFromEnv() = %d, want %d", got, tt.want)
			}
		})
	}
}
