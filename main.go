package main

import (
	"log"
	"net/http"

	"prueba/api/handlers"
	"prueba/api/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Token
	mux.HandleFunc("/token", handlers.PostToken)

	// Personajes (protegido)
	mux.Handle("/personajes", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPersonajes)))

	log.Println("Servidor corriendo en :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
