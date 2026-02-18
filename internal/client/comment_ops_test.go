package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListComments_UsesOrderByDesc(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(emptyAPIResponse())
	}))
	defer server.Close()

	c := newTestClient(server.URL)

	_, err := c.ListComments(context.Background(), 42, 25, nil)
	if err != nil {
		t.Fatalf("ListComments returned unexpected error: %v", err)
	}

	if !strings.Contains(capturedURL, "orderByDesc=CreateDate") {
		t.Errorf("expected URL to contain 'orderByDesc=CreateDate', got: %s", capturedURL)
	}

	// Ensure the old broken form is not present (URL-encoded or plain)
	if strings.Contains(capturedURL, "orderBy=CreateDate+desc") {
		t.Errorf("expected URL NOT to contain 'orderBy=CreateDate+desc', got: %s", capturedURL)
	}
	if strings.Contains(capturedURL, "orderBy=CreateDate%20desc") {
		t.Errorf("expected URL NOT to contain 'orderBy=CreateDate%%20desc', got: %s", capturedURL)
	}
	// Strip the "orderByDesc" occurrences and confirm no plain "orderBy=" remains
	stripped := strings.ReplaceAll(capturedURL, "orderByDesc", "")
	if strings.Contains(stripped, "orderBy=") {
		t.Errorf("expected URL NOT to contain plain 'orderBy=' for ListComments, got: %s", capturedURL)
	}
}
