package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// KVStore wraps the Mattermost plugin KV APIs with namespacing helpers.
type KVStore struct {
	api    plugin.API
	prefix string
}

// NewKVStore constructs a KVStore rooted at the provided namespace prefix.
func NewKVStore(api plugin.API, prefix string) *KVStore {
	if prefix == "" {
		prefix = "ai-suite"
	}

	return &KVStore{
		api:    api,
		prefix: prefix,
	}
}

// SetJSON encodes the given value as JSON and stores it at the provided key.
func (s *KVStore) SetJSON(_ context.Context, key string, value any, expiry time.Duration) error {
	if s == nil || s.api == nil {
		return fmt.Errorf("store: not initialized")
	}

	fullKey, err := s.buildKey(key)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("store: encode value: %w", err)
	}

	if expiry > 0 {
		if appErr := s.api.KVSetWithExpiry(fullKey, payload, int64(expiry.Seconds())); appErr != nil {
			return fmt.Errorf("store: kv set with expiry: %w", appErr)
		}
		return nil
	}

	if appErr := s.api.KVSet(fullKey, payload); appErr != nil {
		return fmt.Errorf("store: kv set: %w", appErr)
	}
	return nil
}

// GetJSON retrieves the value at key and unmarshals into dest.
func (s *KVStore) GetJSON(_ context.Context, key string, dest any) (bool, error) {
	if s == nil || s.api == nil {
		return false, fmt.Errorf("store: not initialized")
	}

	fullKey, err := s.buildKey(key)
	if err != nil {
		return false, err
	}

	data, appErr := s.api.KVGet(fullKey)
	if appErr != nil {
		return false, fmt.Errorf("store: kv get: %w", appErr)
	}

	if len(data) == 0 {
		return false, nil
	}

	if err = json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("store: decode value: %w", err)
	}

	return true, nil
}

// Delete removes the value stored at key.
func (s *KVStore) Delete(_ context.Context, key string) error {
	if s == nil || s.api == nil {
		return fmt.Errorf("store: not initialized")
	}

	fullKey, err := s.buildKey(key)
	if err != nil {
		return err
	}

	if appErr := s.api.KVDelete(fullKey); appErr != nil {
		return fmt.Errorf("store: kv delete: %w", appErr)
	}
	return nil
}

func (s *KVStore) buildKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrInvalidKey
	}

	return fmt.Sprintf("%s:%s", s.prefix, key), nil
}
