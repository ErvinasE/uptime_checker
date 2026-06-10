package api

import (
	"log"
	"net/http"
	"path"
	"strings"
)

// loggerMiddleware logs incoming requests for debugging.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware allows the Next.js frontend to call the API from the browser.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewRouter registers all API routes on a ServeMux.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/websites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.ListWebsites(w, r)
	})
	mux.HandleFunc("/websites/", func(w http.ResponseWriter, r *http.Request) {
		// path.Clean removes double slashes and trailing slashes
		p := path.Clean(r.URL.Path)
		p = strings.Trim(p, "/")
		parts := strings.Split(p, "/")

		// Case: websites/:id/check
		if len(parts) == 3 && parts[0] == "websites" && parts[2] == "check" {
			if r.Method != http.MethodPost {
				writeError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			h.TriggerCheck(w, r)
			return
		}

		// Case: websites/:id/history
		if len(parts) == 3 && parts[0] == "websites" && parts[2] == "history" {
			if r.Method != http.MethodGet {
				writeError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			h.GetWebsiteHistory(w, r)
			return
		}

		// Case: websites/:id
		if len(parts) == 2 && parts[0] == "websites" {
			if r.Method != http.MethodGet {
				writeError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			h.GetWebsite(w, r)
			return
		}

		writeError(w, http.StatusNotFound, "not found")
	})

	return loggerMiddleware(corsMiddleware(mux))
}
