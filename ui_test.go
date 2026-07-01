package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexServesHTMLPage(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(newPacker(challengeSizes, testMaxQuantity)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("content-type: got %q, want text/html", ct)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<form") || !strings.Contains(body, "fetch('/pack'") {
		t.Errorf("page is missing the form or the fetch call to the API")
	}
}
