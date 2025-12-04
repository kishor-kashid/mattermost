package store

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Service exposes high-level helpers for plugin data.
type Service struct {
	kv *KVStore
}

// NewService builds a Service around the provided KV store.
func NewService(kv *KVStore) *Service {
	return &Service{kv: kv}
}

// SaveJSON stores a JSON encodable payload under the given collection/key.
func (s *Service) SaveJSON(ctx context.Context, collection, key string, value any, expiry time.Duration) error {
	if s == nil || s.kv == nil {
		return fmt.Errorf("store: service not initialized")
	}

	return s.kv.SetJSON(ctx, namespacedKey(collection, key), value, expiry)
}

// LoadJSON retrieves a JSON payload into dest. Returns false if not found.
func (s *Service) LoadJSON(ctx context.Context, collection, key string, dest any) (bool, error) {
	if s == nil || s.kv == nil {
		return false, fmt.Errorf("store: service not initialized")
	}

	return s.kv.GetJSON(ctx, namespacedKey(collection, key), dest)
}

// Delete removes the stored value.
func (s *Service) Delete(ctx context.Context, collection, key string) error {
	if s == nil || s.kv == nil {
		return fmt.Errorf("store: service not initialized")
	}

	return s.kv.Delete(ctx, namespacedKey(collection, key))
}

func namespacedKey(parts ...string) string {
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return strings.Join(clean, ":")
}
