package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFile(t *testing.T) (string, func()) {
	t.Helper()

	// Create a temp directory
	dir, err := os.MkdirTemp("", "filestore_test")
	require.NoError(t, err)

	// Create a test file path
	filePath := filepath.Join(dir, "testdb.json")

	// Return cleanup function
	return filePath, func() {
		os.RemoveAll(dir)
	}
}

func TestFileStore_Get(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	// Initialize with test data
	initialData := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}
	writeTestData(t, filePath, initialData)

	store := NewFileStore(filePath).(*FileStore)

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
			url, ok := store.Get(context.Background(), tt.key)
			assert.Equal(t, tt.expectedURL, url)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestFileStore_Set(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	store := NewFileStore(filePath).(*FileStore)

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
		{
			name:        "duplicate key",
			key:         "abc123",
			value:       "http://example.com",
			expectedErr: fmt.Errorf("409 - conflict"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(context.Background(), tt.key, tt.value)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				return
			}

			assert.NoError(t, err)

			// For successful cases, verify the value was actually set
			if tt.expectedErr == nil {
				url, ok := store.Get(context.Background(), tt.key)
				assert.True(t, ok)
				assert.Equal(t, tt.value, url)
			}
		})
	}

	// Test conflict case separately
	t.Run("conflict", func(t *testing.T) {
		key := "conflictKey"
		value := "http://example.com"

		// First set should succeed
		err := store.Set(context.Background(), key, value)
		assert.NoError(t, err)

		// Second set should return conflict
		err = store.Set(context.Background(), key, "http://another.com")
		assert.IsType(t, &appErrors.HTTPError{}, err)
		assert.Equal(t, 409, err.(*appErrors.HTTPError).Code)
	})
}

func TestFileStore_List(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	// Initialize with test data
	initialData := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}
	writeTestData(t, filePath, initialData)

	store := NewFileStore(filePath).(*FileStore)

	result, err := store.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, initialData, result)
}

func TestFileStore_Ping(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	store := NewFileStore(filePath).(*FileStore)
	err := store.Ping(context.Background())
	assert.NoError(t, err)
}

func TestFileStore_SetMultiple(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	store := NewFileStore(filePath).(*FileStore)

	items := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}

	err := store.SetMultiple(context.Background(), items)
	assert.NoError(t, err)

	// Verify all items were set
	for k, v := range items {
		url, ok := store.Get(context.Background(), k)
		assert.True(t, ok)
		assert.Equal(t, v, url)
	}
}

func TestFileStore_Delete(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	// Initialize with test data
	initialData := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}
	writeTestData(t, filePath, initialData)

	store := NewFileStore(filePath).(*FileStore)

	// Delete one item
	err := store.Delete(context.Background(), []string{"abc123"})
	assert.NoError(t, err)

	// Verify deletion
	url, ok := store.Get(context.Background(), "abc123")
	assert.True(t, ok)       // Key still exists
	assert.Equal(t, "", url) // But value is empty

	// Other item should remain unchanged
	url, ok = store.Get(context.Background(), "def456")
	assert.True(t, ok)
	assert.Equal(t, "http://test.org", url)
}

func TestNewFileStore(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	// Write initial test data
	initialData := map[string]string{
		"abc123": "http://example.com",
	}
	writeTestData(t, filePath, initialData)

	store := NewFileStore(filePath).(*FileStore)

	// Verify initialization
	assert.Equal(t, filePath, store.pathToFile)
	assert.Equal(t, initialData, store.s)
}

func TestWriteToFile(t *testing.T) {
	filePath, cleanup := setupTestFile(t)
	defer cleanup()

	// Test writing to file
	err := WriteToFile("abc123", "http://example.com", filePath)
	assert.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	var item itemURL
	err = json.Unmarshal(content, &item)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", item.Hash)
	assert.Equal(t, "http://example.com", item.URL)
}

func writeTestData(t *testing.T, filePath string, data map[string]string) {
	t.Helper()

	file, err := os.Create(filePath)
	require.NoError(t, err)
	defer file.Close()

	for k, v := range data {
		item := itemURL{Hash: k, URL: v}
		jsonData, err := json.Marshal(item)
		require.NoError(t, err)
		_, err = file.WriteString(string(jsonData) + "\n")
		require.NoError(t, err)
	}
}
