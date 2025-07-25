package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test types for testing
type testKey struct {
	ID   int
	Name string
}

func (k testKey) String() string {
	return fmt.Sprintf("%d:%s", k.ID, k.Name)
}

type testValue struct {
	Data string
}

func (v testValue) IsValid() bool {
	return v.Data != ""
}

func TestNewLRUCache(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		capacity int
		wantErr  bool
	}{
		{
			name:     "valid capacity",
			capacity: 100,
			wantErr:  false,
		},
		{
			name:     "zero capacity",
			capacity: 0,
			wantErr:  true,
		},
		{
			name:     "negative capacity",
			capacity: -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache, err := NewLRUCache[testKey, testValue](tt.capacity)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cache)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cache)
				assert.Equal(t, tt.capacity, cache.Capacity())
				assert.Equal(t, 0, cache.Size())
			}
		})
	}
}

func TestLRUCache_GetSet(t *testing.T) {
	t.Parallel()

	cache, err := NewLRUCache[testKey, testValue](10)
	require.NoError(t, err)

	ctx := context.Background()

	key := testKey{
		ID:   1,
		Name: "test",
	}

	value := testValue{
		Data: "test data",
	}

	// Test empty cache
	_, found := cache.Get(ctx, key)
	assert.False(t, found)

	// Test set and get
	err = cache.Set(ctx, key, value)
	assert.NoError(t, err)
	assert.Equal(t, 1, cache.Size())

	retrieved, found := cache.Get(ctx, key)
	assert.True(t, found)
	assert.Equal(t, value, retrieved)
}

func TestLRUCache_Delete(t *testing.T) {
	t.Parallel()

	cache, err := NewLRUCache[testKey, testValue](10)
	require.NoError(t, err)

	ctx := context.Background()

	key := testKey{
		ID:   1,
		Name: "test",
	}

	value := testValue{
		Data: "test data",
	}

	// Set value
	err = cache.Set(ctx, key, value)
	assert.NoError(t, err)
	assert.Equal(t, 1, cache.Size())

	// Delete value
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, 0, cache.Size())

	// Verify value is gone
	_, found := cache.Get(ctx, key)
	assert.False(t, found)
}

func TestLRUCache_Clear(t *testing.T) {
	t.Parallel()

	cache, err := NewLRUCache[testKey, testValue](10)
	require.NoError(t, err)

	ctx := context.Background()

	// Add multiple values
	keys := []testKey{
		{ID: 1, Name: "test1"},
		{ID: 2, Name: "test2"},
		{ID: 3, Name: "test3"},
	}

	for i, key := range keys {
		value := testValue{Data: fmt.Sprintf("data%d", i+1)}
		err := cache.Set(ctx, key, value)
		assert.NoError(t, err)
	}

	assert.Equal(t, 3, cache.Size())

	// Clear cache
	err = cache.Clear(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, cache.Size())

	// Verify all values are gone
	for _, key := range keys {
		_, found := cache.Get(ctx, key)
		assert.False(t, found)
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	t.Parallel()

	cache, err := NewLRUCache[testKey, testValue](2)
	require.NoError(t, err)

	ctx := context.Background()

	// Add 3 values to a cache with capacity 2
	keys := []testKey{
		{ID: 1, Name: "test1"},
		{ID: 2, Name: "test2"},
		{ID: 3, Name: "test3"},
	}

	for i, key := range keys {
		value := testValue{Data: fmt.Sprintf("data%d", i+1)}
		err := cache.Set(ctx, key, value)
		assert.NoError(t, err)
	}

	// Cache should have size 2 (evicted the least recently used)
	assert.Equal(t, 2, cache.Size())

	// First key should be evicted (LRU)
	_, found := cache.Get(ctx, keys[0])
	assert.False(t, found)

	// Last two keys should still be there
	_, found = cache.Get(ctx, keys[1])
	assert.True(t, found)

	_, found = cache.Get(ctx, keys[2])
	assert.True(t, found)
}

func TestCacheKey_String(t *testing.T) {
	t.Parallel()

	key := testKey{
		ID:   123,
		Name: "test",
	}

	assert.Equal(t, "123:test", key.String())
}

func TestCacheValue_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value testValue
		want  bool
	}{
		{
			name:  "valid value",
			value: testValue{Data: "test"},
			want:  true,
		},
		{
			name:  "invalid value",
			value: testValue{Data: ""},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.value.IsValid())
		})
	}
}
