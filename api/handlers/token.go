package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TokenRequest struct {
	Nombre string `json:"nombre"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

var TokenStore = make(map[string]time.Time)

func PostToken(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Nombre == "" {
		http.Error(w, "Debe enviar un nombre válido", http.StatusBadRequest)
		return
	}

	token := uuid.New().String()
	TokenStore[token] = time.Now().Add(10 * time.Minute) // expiración temporal

	resp := TokenResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
