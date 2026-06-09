package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ervinas/server_uptime_checker/internal/models"
)

func TestCheckUp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New()
	result := c.Check(context.Background(), server.URL)

	if result.Status != models.StatusUp {
		t.Fatalf("expected up, got %s", result.Status)
	}
	if result.ResponseTimeMs == nil || *result.ResponseTimeMs < 0 {
		t.Fatal("expected positive response time")
	}
}

func TestCheckDownOnServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := New()
	result := c.Check(context.Background(), server.URL)

	if result.Status != models.StatusDown {
		t.Fatalf("expected down, got %s", result.Status)
	}
	if result.ErrorMessage == nil {
		t.Fatal("expected error message")
	}
}

func TestCheckUpOnForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c := New()
	result := c.Check(context.Background(), server.URL)

	if result.Status != models.StatusUp {
		t.Fatalf("expected up for 403 bot protection, got %s", result.Status)
	}
}

func TestCheckDownOnInvalidURL(t *testing.T) {
	c := New()
	result := c.Check(context.Background(), "http://localhost:1")

	if result.Status != models.StatusDown {
		t.Fatalf("expected down, got %s", result.Status)
	}
}
