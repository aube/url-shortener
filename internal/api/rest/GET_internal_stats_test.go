package restapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStatsGetter struct {
	mock.Mock
}

func (m *MockStatsGetter) GetStats(ctx context.Context) ([]byte, error) {
	args := m.Called(ctx)
	return args.Get(0).([]byte), args.Error(1)
}

func TestHandlerInternalStats(t *testing.T) {

	tests := []struct {
		name           string
		trustedSubnet  string
		xRealIP        string
		mockResponse   []byte
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful stats",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			mockResponse:   []byte(`{"urls":10,"users":5}`),
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "untrusted IP",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "10.0.0.1",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "no trusted subnet",
			trustedSubnet:  "",
			xRealIP:        "192.168.1.100",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "stats error",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override config for test
			originalConfig := config.NewConfig()
			testConfig := originalConfig
			testConfig.TrustedSubnet = tt.trustedSubnet
			config.SetGlobalConfig(testConfig)
			defer func() { config.SetGlobalConfig(originalConfig) }()

			mockGetter := new(MockStatsGetter)
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusBadRequest {
				mockGetter.On("GetStats", mock.Anything).Return(tt.mockResponse, tt.mockError)
			}

			// Replace the real usecase with our mock
			originalGet := usecasesGetStats
			usecasesGetStats = mockGetter.GetStats
			defer func() { usecasesGetStats = originalGet }()

			req := httptest.NewRequest("GET", "/internal/stats", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			rr := httptest.NewRecorder()

			handler := HandlerInternalStats()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockGetter.AssertExpectations(t)
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
