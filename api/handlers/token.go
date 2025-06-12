package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Parametro que se usara para generar el token
type TokenRequest struct {
	Name string `json:"name"`
}

// Token temporal generado
type TokenResponse struct {
	Token string `json:"token"`
}

// Se guarda el token
var TokenStore = make(map[string]time.Time)

// Funcion que retorna el token temporal
func PostToken(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Se requiere enviar un nombre válido", http.StatusBadRequest)
		return
	}

	token := uuid.New().String()
	TokenStore[token] = time.Now().Add(10 * time.Minute) // expiración temporal

	resp := TokenResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
