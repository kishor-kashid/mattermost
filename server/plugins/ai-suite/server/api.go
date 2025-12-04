package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	apiBasePath    = "/api/v1"
	publicBasePath = "/api/v1/public"
)

// APIRouter wires plugin HTTP handlers with middleware.
type APIRouter struct {
	plugin *Plugin
	router *mux.Router
}

// NewAPIRouter configures all HTTP routes for the plugin.
func NewAPIRouter(p *Plugin) *APIRouter {
	r := &APIRouter{
		plugin: p,
		router: mux.NewRouter(),
	}

	r.registerRoutes()
	return r
}

func (r *APIRouter) registerRoutes() {
	public := r.router.PathPrefix(publicBasePath).Subrouter()
	public.HandleFunc("/health", r.handleHealth).Methods(http.MethodGet)

	api := r.router.PathPrefix(apiBasePath).Subrouter()
	api.Use(r.logRequest)
	api.Use(r.requireUser)

	api.HandleFunc("/me", r.handleMe).Methods(http.MethodGet)
}

func (r *APIRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *APIRouter) handleHealth(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func (r *APIRouter) handleMe(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("Mattermost-User-Id")
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": userID,
	})
}

func (r *APIRouter) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.plugin.apiClient.Log.Debug("api request",
			"method", req.Method,
			"path", req.URL.Path,
			"user_id", req.Header.Get("Mattermost-User-Id"),
		)
		next.ServeHTTP(w, req)
	})
}

func (r *APIRouter) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID := strings.TrimSpace(req.Header.Get("Mattermost-User-Id"))
		if userID == "" {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(w, req)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
