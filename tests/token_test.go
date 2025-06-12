package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"prueba/api/handlers"
)

// Test que prueba que el token fue generado con exito
func TestPostTokenSuccess(t *testing.T) {
	body := []byte(`{"name": "Morty"}`)
	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.PostToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Esperado 200 OK, recibido %d", resp.StatusCode)
	}

	var data map[string]string
	bodyResp, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyResp, &data); err != nil {
		t.Fatal("No se pudo decodificar respuesta JSON")
	}

	if _, ok := data["token"]; !ok {
		t.Fatal("Respuesta no contiene 'token'")
	}
}

// Test que prueba que el token necesita el parametro name
func TestPostTokenMissingName(t *testing.T) {
	body := []byte(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.PostToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Esperado 400 BadRequest, recibido %d", resp.StatusCode)
	}
}
