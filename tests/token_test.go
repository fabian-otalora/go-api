package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"prueba/api/handlers"
)

// Test que prueba generacion de token
func TestGenerateToken(t *testing.T) {
	// Simular variables del entorno
	os.Setenv("TOKEN_EXPIRATION_MINUTES", "15")

	body := `{"name": "Morty"}`
	req := httptest.NewRequest("POST", "/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handlers.GenerateToken(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperado 200 OK, obtuve %d", rec.Code)
	}

	var resp struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}

	err := json.NewDecoder(rec.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("error al decodificar respuesta: %v", err)
	}

	if resp.Token == "" {
		t.Error("token vacío en respuesta")
	}
}
