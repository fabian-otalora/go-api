package handlers

import (
	"encoding/json"
	"net/http"
	"os"
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
	apiURL := os.Getenv("RICK_AND_MORTY_API")
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "No se pudo obtener personajes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error al obtener personajes", resp.StatusCode)
		return
	}

	var data struct {
		Results []map[string]interface{} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Error al parsear JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Results)
}
