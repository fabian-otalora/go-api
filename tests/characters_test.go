package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"prueba/api/handlers"
	"prueba/api/middleware"
)

// Test que prueba el limite de intentos
func TestTokenAttemptLimit(t *testing.T) {
	os.Setenv("TOKEN_ATTEMPTS_LIMIT", "5")

	token := "test-token"
	exp := time.Now().Add(10 * time.Minute)
	handlers.TokenStore[token] = exp

	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Hacer 5 peticiones válidas
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/characters", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Esperado 200 OK en intento %d, pero obtuve %d", i+1, rec.Code)
		}
	}

	// 6ta petición debe fallar
	req := httptest.NewRequest("GET", "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Esperado 429 Too Many Requests en el intento 6, pero obtuve %d", rec.Code)
	}
}

// Test que prueba la expiracion del token
func TestTokenExpired(t *testing.T) {
	token := "expired-token"
	handlers.TokenStore[token] = time.Now().Add(-5 * time.Minute) // expirado

	req := httptest.NewRequest("GET", "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 Unauthorized por token expirado, obtuve %d", rec.Code)
	}
}
