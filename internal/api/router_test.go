package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ervinas/server_uptime_checker/internal/checker"
	"github.com/ervinas/server_uptime_checker/internal/store"
)

func TestRouter_TriggerCheck_TrailingSlash(t *testing.T) {
	// We need a store and checker, but for routing test we can try to use nil or minimal mocks if possible.
	// Since NewHandler takes *store.Store and *checker.Checker, and they are structs, we might need real ones
	// or they might panic. Let's use real ones with a nil DB (might still panic) or a simple mock.
	
	// For simplicity, let's just test the logic in a more isolated way if possible, 
	// or just provide what it needs.
	s := store.New(nil)
	c := checker.New()
	h := NewHandler(s, c)
	router := NewRouter(h)

	// Test POST /websites/1/check
	req := httptest.NewRequest(http.MethodPost, "/websites/1/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// We expect 500 (failed to get website because store is nil) or 400, but NOT 405.
	if w.Code == http.StatusMethodNotAllowed {
		t.Errorf("expected NOT 405 for /websites/1/check, got %d", w.Code)
	}

	// Test POST /websites/1/check/
	req = httptest.NewRequest(http.MethodPost, "/websites/1/check/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusMethodNotAllowed {
		t.Errorf("expected NOT 405 for /websites/1/check/, got %d", w.Code)
	}
}
