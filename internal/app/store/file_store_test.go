package store

import (
	"context"
	"os"
	"testing"

	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	// Setup test directory
	tempDir := t.TempDir()
	ctx := context.WithValue(context.Background(), ctxkeys.UserIDKey, "test-user")

	t.Run("NewFileStore creates files", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		assert.FileExists(t, store.urlsFile)
		assert.FileExists(t, store.usersFile)
	})

	t.Run("Set and Get", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "abc123", "https://example.com")
		require.NoError(t, err)

		// Test Get
		val, ok := store.GetByUser(ctx, "abc123")
		assert.True(t, ok)
		assert.Equal(t, "https://example.com", val)

		// Test Get with wrong user
		wrongUserCtx := context.WithValue(context.Background(), ctxkeys.UserIDKey, "wrong-user")
		val, ok = store.GetByUser(wrongUserCtx, "abc123")
		assert.False(t, ok)
		assert.Empty(t, val)
	})

	t.Run("Set duplicate key", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "abc123", "https://example.com")
		require.NoError(t, err)

		err = store.Set(ctx, "abc123", "https://example.org")
		require.Error(t, err)
		assert.Equal(t, "409 - conflict", err.Error())
	})

	t.Run("List", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "key1", "https://example.com/1")
		require.NoError(t, err)
		err = store.Set(ctx, "key2", "https://example.com/2")
		require.NoError(t, err)

		// Test with different user
		otherUserCtx := context.WithValue(context.Background(), ctxkeys.UserIDKey, "other-user")
		err = store.Set(otherUserCtx, "key3", "https://example.com/3")
		require.NoError(t, err)

		items, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, "https://example.com/1", items["key1"])
		assert.Equal(t, "https://example.com/2", items["key2"])
	})

	t.Run("SetMultiple", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		items := map[string]string{
			"key1": "https://example.com/1",
			"key2": "https://example.com/2",
		}

		err := store.SetMultiple(ctx, items)
		require.NoError(t, err)

		val, ok := store.Get(ctx, "key1")
		assert.True(t, ok)
		assert.Equal(t, "https://example.com/1", val)

		val, ok = store.Get(ctx, "key2")
		assert.True(t, ok)
		assert.Equal(t, "https://example.com/2", val)
	})

	t.Run("Delete", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "key1", "https://example.com/1")
		require.NoError(t, err)
		err = store.Set(ctx, "key2", "https://example.com/2")
		require.NoError(t, err)

		err = store.Delete(ctx, []string{"key1", "key2"})
		require.NoError(t, err)

		val, ok := store.Get(ctx, "key1")
		assert.False(t, ok)
		assert.Empty(t, val)

		val, ok = store.Get(ctx, "key2")
		assert.False(t, ok)
		assert.Empty(t, val)
	})

	t.Run("Stats", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "key1", "https://example.com/1")
		require.NoError(t, err)
		err = store.Set(ctx, "key2", "https://example.com/2")
		require.NoError(t, err)

		// Different user
		otherUserCtx := context.WithValue(context.Background(), ctxkeys.UserIDKey, "other-user")
		err = store.Set(otherUserCtx, "key3", "https://example.com/3")
		require.NoError(t, err)

		urlCount, userCount, err := store.Stats(ctx)
		require.NoError(t, err)
		assert.Equal(t, 3, urlCount)
		assert.Equal(t, 2, userCount)
	})

	t.Run("Invalid inputs", func(t *testing.T) {
		store := NewFileStore(tempDir).(*FileStore)
		defer os.RemoveAll(tempDir)

		err := store.Set(ctx, "", "https://example.com")
		require.Error(t, err)
		assert.Equal(t, "invalid input", err.Error())

		err = store.Set(ctx, "key", "")
		require.Error(t, err)
		assert.Equal(t, "invalid input", err.Error())
	})
}
