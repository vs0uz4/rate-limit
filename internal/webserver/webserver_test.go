package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/rate-limit/config"
	"github.com/vs0uz4/rate-limit/internal/mock"
)

func TestPingEndpoint(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Unexpected status code to /ping")
	assert.Equal(t, "pong", w.Body.String(), "Unexpected response to /ping")
}

func TestStartServer(t *testing.T) {
	mockLimiter := &mock.MockRateLimiter{Allowed: true}

	cfg := &config.Config{
		WebServerPort: "9090",
	}

	go func() {
		err := Start(cfg, mockLimiter)
		assert.NoError(t, err, "Unexpected error when start webserver")
	}()

	resp, err := http.Get("http://localhost:9090/ping")
	assert.NoError(t, err, "Error when accessing the server")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")
}
