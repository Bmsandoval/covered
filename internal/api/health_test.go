package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bmsandoval/covered/internal/app"
	"github.com/Bmsandoval/covered/internal/config"
)

func TestHealth(t *testing.T) {
	srv := NewServer(config.Config{AppPort: 8080}, app.Deps{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" || body.Phase != "prototype" {
		t.Fatalf("body = %+v", body)
	}
}
