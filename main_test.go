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

func TestPackReturns200(t *testing.T) {
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, newPackRequest(t, 1))

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestPackForQuantityOf1ReturnsSingle250Pack(t *testing.T) {
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, newPackRequest(t, 1))

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

	handler().ServeHTTP(rec, newPackRequest(t, 251))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	got := decodePackResponse(t, rec.Body)

	want := packResponse{Packs: map[int]int{500: 1}}
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
