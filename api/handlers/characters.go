package handlers

import (
	"encoding/json"
	"net/http"
)

type RickAndMortyResponse struct {
	Results []Character `json:"results"`
}

// Data que se recibe de la API
type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Gender  string `json:"gender"`
	Species string `json:"species"`
}

// Obtener lista de personajes
func GetCharacters(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://rickandmortyapi.com/api/character")
	if err != nil {
		http.Error(w, "Error ", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data RickAndMortyResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Error al procesar datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Results)
}
