package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var challengeSizes = []int{250, 500, 1000, 2000, 5000}

const testMaxQuantity = 100_000

func TestPackReturns200(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(newPacker(challengeSizes, testMaxQuantity)).ServeHTTP(rec, newPackRequest(t, 1))

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content-type: got %q, want %q", ct, "application/json")
	}
}

func TestPackRejectsNonPostMethod(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/pack", nil)

	handler(newPacker(challengeSizes, testMaxQuantity)).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if allow := rec.Header().Get("Allow"); allow != http.MethodPost {
		t.Errorf("allow: got %q, want %q", allow, http.MethodPost)
	}
}

func TestPackCalculatesPacks(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []int
		quantity int
		want     map[int]int
	}{
		{"quantity 1 rounds up to one 250", challengeSizes, 1, map[int]int{250: 1}},
		{"quantity 250 fits exactly in one 250", challengeSizes, 250, map[int]int{250: 1}},
		{"quantity 251 rounds up to one 500", challengeSizes, 251, map[int]int{500: 1}},
		{"quantity 501 combines a 500 and a 250", challengeSizes, 501, map[int]int{500: 1, 250: 1}},
		{"quantity 751 rounds up to one 1000", challengeSizes, 751, map[int]int{1000: 1}},
		{"quantity 1751 rounds up to one 2000", challengeSizes, 1751, map[int]int{2000: 1}},
		{"quantity 4751 rounds up to one 5000", challengeSizes, 4751, map[int]int{5000: 1}},
		{"quantity 12001 combines two 5000s, a 2000 and a 250", challengeSizes, 12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
		{"non-multiple sizes 4 and 5 for quantity 7 use two 4s", []int{4, 5}, 7, map[int]int{4: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			handler(newPacker(tt.sizes, testMaxQuantity)).ServeHTTP(rec, newPackRequest(t, tt.quantity))

			if rec.Code != http.StatusOK {
				t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
			}

			got := decodePackResponse(t, rec.Body)

			want := packResponse{Packs: tt.want}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
			}
		})
	}
}

func TestPackUsesConfiguredPackSizes(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(newPacker([]int{100, 300}, testMaxQuantity)).ServeHTTP(rec, newPackRequest(t, 100))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{100: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackRejectsQuantityAboveMax(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(newPacker(challengeSizes, 1000)).ServeHTTP(rec, newPackRequest(t, 1001))

	assertBadRequest(t, rec, "quantity must not exceed 1000")
}

func TestPackRejectsQuantityBelowOne(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(newPacker(challengeSizes, testMaxQuantity)).ServeHTTP(rec, newPackRequest(t, 0))

	assertBadRequest(t, rec, "quantity must be at least 1")
}

func TestPackRejectsOversizedBody(t *testing.T) {
	var body bytes.Buffer
	body.WriteString(`{"quantity":251,"padding":"`)
	body.Write(bytes.Repeat([]byte("a"), 2<<20))
	body.WriteString(`"}`)

	req := httptest.NewRequest(http.MethodPost, "/pack", &body)
	rec := httptest.NewRecorder()

	handler(newPacker(challengeSizes, testMaxQuantity)).ServeHTTP(rec, req)

	assertBadRequest(t, rec, "invalid request body")
}

func newPackRequest(t *testing.T, quantity int) *http.Request {
	t.Helper()
	body, err := json.Marshal(struct {
		Quantity int `json:"quantity"`
	}{Quantity: quantity})
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	return httptest.NewRequest(http.MethodPost, "/pack", bytes.NewReader(body))
}

func decodePackResponse(t *testing.T, body io.Reader) packResponse {
	t.Helper()
	var got packResponse
	if err := json.NewDecoder(body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return got
}

func assertBadRequest(t *testing.T, rec *httptest.ResponseRecorder, detail string) {
	t.Helper()

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("content-type: got %q, want %q", ct, "application/problem+json")
	}

	var prob problemDetail
	if err := json.NewDecoder(rec.Body).Decode(&prob); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if prob.Status != http.StatusBadRequest {
		t.Errorf("status field: got %d, want %d", prob.Status, http.StatusBadRequest)
	}
	if prob.Title != http.StatusText(http.StatusBadRequest) {
		t.Errorf("title: got %q, want %q", prob.Title, http.StatusText(http.StatusBadRequest))
	}
	if prob.Detail != detail {
		t.Errorf("detail: got %q, want %q", prob.Detail, detail)
	}
}
