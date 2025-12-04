package summarizer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/store"
)

const summaryCollection = "summaries"

// Cache persists summaries in the plugin KV store to reduce LLM calls.
type Cache struct {
	store *store.Service
	ttl   time.Duration
}

// NewCache creates a cache with the provided TTL (defaults to 24h).
func NewCache(store *store.Service, ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &Cache{
		store: store,
		ttl:   ttl,
	}
}

// Get retrieves a cached summary by key.
func (c *Cache) Get(ctx context.Context, key string) (*Response, bool, error) {
	if c == nil || c.store == nil || key == "" {
		return nil, false, nil
	}

	var entry Response
	found, err := c.store.LoadJSON(ctx, summaryCollection, key, &entry)
	if err != nil || !found {
		return nil, false, err
	}
	entry.Cached = true
	return &entry, true, nil
}

// Set stores the given summary under key.
func (c *Cache) Set(ctx context.Context, key string, summary *Response) error {
	if c == nil || c.store == nil || key == "" || summary == nil {
		return nil
	}

	copy := *summary
	copy.Cached = false
	return c.store.SaveJSON(ctx, summaryCollection, key, copy, c.ttl)
}

// BuildCacheKey joins and hashes arbitrary parts into a deterministic key.
func BuildCacheKey(parts ...string) string {
	joined := strings.Join(parts, "|")
	sum := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(sum[:])
}
