package store

import (
	"context"
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/require"
)

func TestServiceSaveAndLoadJSON(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVSet", "ai-suite:test:key", []byte(`{"value":"demo"}`)).Return(nil)
	api.On("KVGet", "ai-suite:test:key").Return([]byte(`{"value":"demo"}`), nil)

	svc := NewService(NewKVStore(api, "ai-suite"))

	type sample struct {
		Value string `json:"value"`
	}

	err := svc.SaveJSON(context.Background(), "test", "key", sample{Value: "demo"}, 0)
	require.NoError(t, err)

	var dest sample
	found, err := svc.LoadJSON(context.Background(), "test", "key", &dest)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "demo", dest.Value)
}

func TestServiceDelete(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVDelete", "ai-suite:test:key").Return(nil)

	svc := NewService(NewKVStore(api, "ai-suite"))
	err := svc.Delete(context.Background(), "test", "key")
	require.NoError(t, err)
}

func TestKVSetJSONWithExpiry(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVSetWithExpiry", "ai-suite:exp:key", []byte(`{"value":"demo"}`), int64(5)).Return(nil)

	kv := NewKVStore(api, "ai-suite")
	err := kv.SetJSON(context.Background(), "exp:key", map[string]string{"value": "demo"}, 5*time.Second)
	require.NoError(t, err)
}
