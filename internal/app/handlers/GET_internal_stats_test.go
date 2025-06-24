package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorageStats is a mock implementation of StorageStats for testing
type MockStorageStats struct {
	mock.Mock
}

func (m *MockStorageStats) Stats(ctx context.Context) (int, int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Int(1), args.Error(2)
}

func TestHandlerInternalStats(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name              string
		trustedSubnet     string
		xRealIP           string
		mockUrls          int
		mockUsers         int
		mockError         error
		expectedStatus    int
		expectedResponse  string
		expectErrorLogged bool
	}{
		{
			name:             "successful request from trusted IP",
			trustedSubnet:    "192.168.1.0/24",
			xRealIP:          "192.168.1.100",
			mockUrls:         42,
			mockUsers:        10,
			mockError:        nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"urls":42,"users":10}`,
		},
		{
			name:             "untrusted IP",
			trustedSubnet:    "192.168.1.0/24",
			xRealIP:          "10.0.0.1",
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "wrong IP\n",
		},
		{
			name:             "no trusted subnet configured",
			trustedSubnet:    "",
			xRealIP:          "192.168.1.100",
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "wrong IP\n",
		},
		{
			name:              "storage error",
			trustedSubnet:     "192.168.1.0/24",
			xRealIP:           "192.168.1.100",
			mockError:         assert.AnError,
			expectedStatus:    http.StatusBadRequest,
			expectedResponse:  "Error on stats reuest\n",
			expectErrorLogged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock storage
			mockStorage := new(MockStorageStats)
			mockStorage.On("Stats", mock.Anything).
				Return(tt.mockUrls, tt.mockUsers, tt.mockError)

			// Override config for test
			originalConfig := config.NewConfig()
			defer func() { config.SetGlobalConfig(originalConfig) }()

			testConfig := config.NewConfig()
			testConfig.TrustedSubnet = tt.trustedSubnet
			config.SetGlobalConfig(testConfig)

			// Create request with X-Real-IP header
			req := httptest.NewRequest("GET", "/internal/stats", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call handler
			handler := HandlerInternalStats(mockStorage)
			handler(rec, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedResponse, rec.Body.String())

			// Verify mock expectations
			if tt.mockError == nil && tt.expectedStatus == http.StatusOK {
				mockStorage.AssertCalled(t, "Stats", req.Context())
			}
		})
	}
}

func TestIsIPInCIDR(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
	}{
		{"IP in CIDR", "192.168.1.100", "192.168.1.0/24", true},
		{"IP not in CIDR", "10.0.0.1", "192.168.1.0/24", false},
		{"Invalid IP", "invalid", "192.168.1.0/24", false},
		{"Invalid CIDR", "192.168.1.100", "invalid", false},
		{"IPv6 in CIDR", "2001:db8::1", "2001:db8::/32", true},
		{"IPv6 not in CIDR", "2001:db8::1", "2002:db8::/32", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPInCIDR(tt.ip, tt.cidr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
