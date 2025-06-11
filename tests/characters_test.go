package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"prueba/api/handlers"
	"prueba/api/middleware"
)

func TestGetPersonajesSuccess(t *testing.T) {
	// Crear token manualmente y registrarlo
	token := "test-token"
	handlers.TokenStore[token] = time.Now().Add(10 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/personajes", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPersonajes))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado 200 OK, recibido %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatal("no se pudo decodificar JSON de personajes")
	}
}

func TestGetPersonajesUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/personajes", nil)
	// sin token
	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPersonajes))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperado 401 Unauthorized, recibido %d", resp.StatusCode)
	}
}

func TestGetPersonajesExpiredToken(t *testing.T) {
	token := "expired-token"
	// Token con tiempo de expiración en el pasado
	handlers.TokenStore[token] = time.Now().Add(-1 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/personajes", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPersonajes))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperado 401 Unauthorized por token expirado, recibido %d", resp.StatusCode)
	}
}

func TestGetPersonajesTooManyRequests(t *testing.T) {
	token := "limited-token"
	handlers.TokenStore[token] = time.Now().Add(10 * time.Minute) // válido

	req := httptest.NewRequest(http.MethodGet, "/personajes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPersonajes))

	var lastStatus int
	for i := 1; i <= 6; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		lastStatus = w.Result().StatusCode
		w.Result().Body.Close()
	}

	if lastStatus != http.StatusTooManyRequests {
		t.Fatalf("esperado 429 Too Many Requests en el intento 6, recibido %d", lastStatus)
	}
}
