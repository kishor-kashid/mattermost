package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChatCompletionSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Fatalf("expected authorization header")
		}
		_ = json.NewEncoder(w).Encode(ChatCompletionResponse{
			ID: "test",
			Choices: []ChatChoice{
				{
					Message: Message{
						Role:    RoleAssistant,
						Content: "ok",
					},
				},
			},
		})
	}))
	defer ts.Close()

	client, err := NewClient(ClientConfig{
		APIKey:               "test",
		Model:                "gpt-4",
		BaseURL:              ts.URL,
		MaxRequestsPerMinute: 10,
		RequestTimeout:       time.Second,
	})
	if err != nil {
		t.Fatalf("expected client, got error %v", err)
	}

	resp, err := client.ChatCompletion(context.Background(), ChatCompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "ping"}},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if resp.ID != "test" {
		t.Fatalf("unexpected response %+v", resp)
	}
}

func TestChatCompletionError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: ErrorDetails{
				Message: "bad request",
				Code:    "invalid_request",
			},
		})
	}))
	defer ts.Close()

	client, err := NewClient(ClientConfig{
		APIKey:         "test",
		Model:          "gpt-4",
		BaseURL:        ts.URL,
		RequestTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("expected client, got error %v", err)
	}

	_, err = client.ChatCompletion(context.Background(), ChatCompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "ping"}},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
