package store

import (
	"context"
	"fmt"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	// Test data
	store.s["abc123"] = "http://example.com"

	tests := []struct {
		name        string
		key         string
		expectedURL string
		expectedOk  bool
	}{
		{
			name:        "existing key",
			key:         "abc123",
			expectedURL: "http://example.com",
			expectedOk:  true,
		},
		{
			name:        "non-existent key",
			key:         "nonexistent",
			expectedURL: "",
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, ok := store.Get(ctx, tt.key)
			assert.Equal(t, tt.expectedURL, url)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestMemoryStore_Set(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	tests := []struct {
		name        string
		key         string
		value       string
		expectedErr error
	}{
		{
			name:        "successful set",
			key:         "abc123",
			value:       "http://example.com",
			expectedErr: nil,
		},
		{
			name:        "empty key",
			key:         "",
			value:       "http://example.com",
			expectedErr: fmt.Errorf("invalid input"),
		},
		{
			name:        "empty value",
			key:         "abc123",
			value:       "",
			expectedErr: fmt.Errorf("invalid input"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(ctx, tt.key, tt.value)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				return
			}

			assert.NoError(t, err)

			// Verify the value was actually set
			url, ok := store.Get(ctx, tt.key)
			assert.True(t, ok)
			assert.Equal(t, tt.value, url)
		})
	}

	// Test conflict case separately
	t.Run("conflict", func(t *testing.T) {
		key := "conflictKey"
		value := "http://example.com"

		// First set should succeed
		err := store.Set(ctx, key, value)
		assert.NoError(t, err)

		// Second set should return conflict
		err = store.Set(ctx, key, "http://another.com")
		assert.IsType(t, &appErrors.HTTPError{}, err)
		assert.Equal(t, 409, err.(*appErrors.HTTPError).Code)
	})
}

func TestMemoryStore_List(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	// Test data
	expectedData := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}
	store.s = expectedData

	result, err := store.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestMemoryStore_Ping(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	err := store.Ping(ctx)
	assert.NoError(t, err)
}

func TestMemoryStore_SetMultiple(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	items := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}

	err := store.SetMultiple(ctx, items)
	assert.NoError(t, err)

	// Verify all items were set
	for k, v := range items {
		url, ok := store.Get(ctx, k)
		assert.True(t, ok)
		assert.Equal(t, v, url)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemStore().(*MemoryStore)
	ctx := context.Background()

	// Initialize with test data
	store.s = map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}

	// Delete one item
	err := store.Delete(ctx, []string{"abc123"})
	assert.NoError(t, err)

	// Verify deletion
	url, ok := store.Get(ctx, "abc123")
	assert.True(t, ok)       // Key still exists
	assert.Equal(t, "", url) // But value is empty

	// Other item should remain unchanged
	url, ok = store.Get(ctx, "def456")
	assert.True(t, ok)
	assert.Equal(t, "http://test.org", url)
}

func TestNewMemStore(t *testing.T) {
	store := NewMemStore().(*MemoryStore)

	// Verify initialization
	assert.NotNil(t, store.s)
	assert.Empty(t, store.s)
}
