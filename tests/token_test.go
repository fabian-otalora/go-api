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

func TestPostTokenSuccess(t *testing.T) {
	body := []byte(`{"nombre": "Morty"}`)
	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.PostToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado 200 OK, recibido %d", resp.StatusCode)
	}

	var data map[string]string
	bodyResp, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyResp, &data); err != nil {
		t.Fatal("no se pudo decodificar respuesta JSON")
	}

	if _, ok := data["token"]; !ok {
		t.Fatal("respuesta no contiene 'token'")
	}
}

func TestPostTokenMissingName(t *testing.T) {
	body := []byte(`{"nombre": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.PostToken(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperado 400 BadRequest, recibido %d", resp.StatusCode)
	}
}
