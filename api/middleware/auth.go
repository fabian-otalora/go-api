package middleware

import (
	"net/http"
	"strings"
	"time"

	"prueba/api/handlers"

	"github.com/patrickmn/go-cache"
)

// Expiracion de tokens
var tokenAttempts = cache.New(10*time.Minute, 15*time.Minute)

// Intentos maximos permitidos del token
const MaxAttempts = 5

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Falta token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		// Validar si el token existe y no está expirado
		exp, ok := handlers.TokenStore[token]
		if !ok || time.Now().After(exp) {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Limitar intentos
		attempts, found := tokenAttempts.Get(token)
		if found && attempts.(int) >= MaxAttempts {
			http.Error(w, "Demasiados intentos, genere un nuevo token", http.StatusTooManyRequests)
			return
		}
		if !found {
			tokenAttempts.Set(token, 1, cache.DefaultExpiration)
		} else {
			tokenAttempts.Set(token, attempts.(int)+1, cache.DefaultExpiration)
		}

		next.ServeHTTP(w, r)
	})
}
