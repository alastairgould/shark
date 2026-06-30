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

func TestPackReturns200(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 1))

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestPackForQuantityOf1ReturnsSingle250Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 1))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{250: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf250ReturnsSingle250Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 250))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{250: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf251ReturnsSingle500Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 251))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{500: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf501Returns500And250Packs(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 501))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{500: 1, 250: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf751ReturnsSingle1000Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 751))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{1000: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf1751ReturnsSingle2000Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 1751))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{2000: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf4751ReturnsSingle5000Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 4751))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{5000: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackForQuantityOf12001Returns5000sAnd2000And250(t *testing.T) {
	rec := httptest.NewRecorder()

	handler(challengeSizes).ServeHTTP(rec, newPackRequest(t, 12001))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{5000: 2, 2000: 1, 250: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

func TestPackUsesConfiguredPackSizes(t *testing.T) {
	rec := httptest.NewRecorder()

	handler([]int{100, 300}).ServeHTTP(rec, newPackRequest(t, 100))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{100: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("packs: got %v, want %v", got.Packs, want.Packs)
	}
}

// newPackRequest builds a POST /pack request with the given quantity as its JSON body.
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

// decodePackResponse reads the JSON body into a packResponse.
func decodePackResponse(t *testing.T, body io.Reader) packResponse {
	t.Helper()
	var got packResponse
	if err := json.NewDecoder(body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return got
}
