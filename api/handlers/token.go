package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Parametro que se usara para generar el token
type TokenRequest struct {
	Name string `json:"name"`
}

// Token temporal generado
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Se guarda el token
var TokenStore = make(map[string]time.Time)

// Funcion que retorna el token temporal
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar parámetro obligatorio
	if req.Name == "" {
		http.Error(w, "El campo 'name' es obligatorio", http.StatusBadRequest)
		return
	}

	// Leer duración desde .env
	expMin := 10 // valor por defecto
	if env := os.Getenv("TOKEN_EXPIRATION_MINUTES"); env != "" {
		if parsed, err := strconv.Atoi(env); err == nil {
			expMin = parsed
		}
	}

	// Generar token
	token := uuid.New().String()
	expiration := time.Now().Add(time.Duration(expMin) * time.Minute)

	// Guardar token con expiración
	TokenStore[token] = expiration
	log.Printf("Token generado para '%s': %s (expira en %v)\n", req.Name, token, expiration)

	// Devolver respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		Token:     token,
		ExpiresAt: expiration,
	})
}
