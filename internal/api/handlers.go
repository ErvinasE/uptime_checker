package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ervinas/server_uptime_checker/internal/store"
)

// Handler serves HTTP API endpoints.
type Handler struct {
	store *store.Store
}

// NewHandler creates an API handler backed by the store.
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListWebsites(w http.ResponseWriter, r *http.Request) {
	websites, err := h.store.ListWebsites()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list websites")
		return
	}
	writeJSON(w, http.StatusOK, websites)
}

func (h *Handler) GetWebsite(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/websites/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid website id")
		return
	}

	website, err := h.store.GetWebsite(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "website not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get website")
		return
	}
	writeJSON(w, http.StatusOK, website)
}

func (h *Handler) GetWebsiteHistory(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/websites/")
	path = strings.TrimSuffix(path, "/history")
	id, err := strconv.ParseInt(strings.Trim(path, "/"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid website id")
		return
	}

	hours := 24
	if raw := r.URL.Query().Get("hours"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			writeError(w, http.StatusBadRequest, "invalid hours parameter")
			return
		}
		hours = parsed
	}

	if _, err := h.store.GetWebsite(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "website not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get website")
		return
	}

	checks, err := h.store.ListChecks(id, hours)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list checks")
		return
	}
	writeJSON(w, http.StatusOK, checks)
}

func parseID(path, prefix string) (int64, error) {
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.Trim(raw, "/")
	return strconv.ParseInt(raw, 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
