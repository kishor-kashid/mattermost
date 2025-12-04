package summarizer

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// HTTPHandler exposes summarization endpoints over HTTP.
type HTTPHandler struct {
	serviceFn func() *Service
}

// NewHTTPHandler wires the handler to the provided service factory.
func NewHTTPHandler(serviceFn func() *Service) *HTTPHandler {
	return &HTTPHandler{serviceFn: serviceFn}
}

// HandleSummarize executes POST /summarize requests.
func (h *HTTPHandler) HandleSummarize(w http.ResponseWriter, r *http.Request) {
	service := h.service()
	if service == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "summarization unavailable")
		return
	}

	userID := strings.TrimSpace(r.Header.Get("Mattermost-User-Id"))
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var payload summarizePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req := Request{
		Type:        Type(strings.ToLower(payload.Type)),
		UserID:      userID,
		ChannelID:   strings.TrimSpace(payload.ChannelID),
		RootPostID:  strings.TrimSpace(payload.RootPostID),
		PostID:      strings.TrimSpace(payload.PostID),
		TimeRange:   strings.TrimSpace(payload.TimeRange),
		SinceMillis: payload.Since,
		UntilMillis: payload.Until,
		Force:       payload.Force,
	}
	if req.Type == "" {
		req.Type = TypeThread
	}

	summary, err := service.Summarize(r.Context(), req)
	if err != nil {
		writeJSONError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func (h *HTTPHandler) service() *Service {
	if h == nil || h.serviceFn == nil {
		return nil
	}
	return h.serviceFn()
}

type summarizePayload struct {
	Type       string `json:"type"`
	ChannelID  string `json:"channel_id"`
	RootPostID string `json:"root_post_id"`
	PostID     string `json:"post_id"`
	TimeRange  string `json:"time_range"`
	Since      int64  `json:"since"`
	Until      int64  `json:"until"`
	Force      bool   `json:"force"`
}

func statusForError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrUnauthorized):
		return http.StatusForbidden
	case errors.Is(err, ErrDisabled), errors.Is(err, ErrClientUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrEmptyConversation):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
