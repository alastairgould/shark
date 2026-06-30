package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPackReturns200(t *testing.T) {
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, newPackRequest(t, 1))

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
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
