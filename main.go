package main

import (
	"log"
	"net/http"
	"os"

	"prueba/api/handlers"
	"prueba/api/middleware"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar .env, usando variables por defecto o de entorno")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Token
	mux.HandleFunc("/token", handlers.GenerateToken)

	// Personajes (se necesita el token para poder usar)
	mux.Handle("/characters", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetCharacters)))

	log.Println("Servidor corriendo en el puerto", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
