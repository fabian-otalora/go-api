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

// Test para probar funcionamiento del servicio obtención de usuarios
func TestGetCharactersSuccess(t *testing.T) {
	// Crear token manualmente y registrarlo
	token := "test-token"
	handlers.TokenStore[token] = time.Now().Add(10 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCharacters))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Se esperaba un 200 OK, recibido %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatal("No se pudo obtener JSON de personajes")
	}
}

// Test para probar funcionamiento del servicio cuando no esta autorizado el usuario
func TestGetCharactersUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	// sin token
	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCharacters))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Se espera 401 Unauthorized, recibido %d", resp.StatusCode)
	}
}

// Test para probar la expiracion del token
func TestGetPersonajesExpiredToken(t *testing.T) {
	token := "expired-token"
	// Token con tiempo de expiración en el pasado
	handlers.TokenStore[token] = time.Now().Add(-1 * time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCharacters))
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Se espera 401 Unauthorized por token expirado, recibido %d", resp.StatusCode)
	}
}

// Test para probar los maximos intentos permitidos del token , los cuales son 5
func TestGetPersonajesTooManyRequests(t *testing.T) {
	token := "limited-token"
	handlers.TokenStore[token] = time.Now().Add(10 * time.Minute) // válido

	req := httptest.NewRequest(http.MethodGet, "/characters", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	handler := middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCharacters))

	var lastStatus int
	for i := 1; i <= 6; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		lastStatus = w.Result().StatusCode
		w.Result().Body.Close()
	}

	if lastStatus != http.StatusTooManyRequests {
		t.Fatalf("Se espera 429 Too Many Requests en el intento 6, recibido %d", lastStatus)
	}
}
