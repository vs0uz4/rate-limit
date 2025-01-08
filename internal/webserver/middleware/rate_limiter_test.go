package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/rate-limit/internal/mock"
)

func TestExtractIPWithoutPort(t *testing.T) {
	ip := extractIP("192.168.1.1")
	assert.Equal(t, "192.168.1.1", ip, "Expected IP without port to be returned as is")
}

func TestExtractIPIPv6ToIPv4(t *testing.T) {
	ip := extractIP("[::1]:8080")
	assert.Equal(t, "127.0.0.1", ip, "Expected IPv6 loopback to be converted to IPv4")
}

func TestRateLimiterMiddlewareTokenValid(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Allowed: true}
	middleware := RateLimiterMiddleware(mockLimiter)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "valid_token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestRateLimiterMiddlewareTokenInvalid(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Allowed: true}
	middleware := RateLimiterMiddleware(mockLimiter)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "invalid_token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestRateLimiterMiddlewareNoToken(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Allowed: true}
	middleware := RateLimiterMiddleware(mockLimiter)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestRateLimiterMiddlewareLimitExceeded(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Allowed: false}
	middleware := RateLimiterMiddleware(mockLimiter)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "You have reached the maximum number of requests or actions allowed within a certain time frame\n", rec.Body.String())
}

func TestRateLimiterMiddlewareError(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Err: errors.New("rate limiter error")}
	middleware := RateLimiterMiddleware(mockLimiter)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "Internal Server Error\n", rec.Body.String())
}
