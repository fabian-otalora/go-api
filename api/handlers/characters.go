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
	if apiURL == "" {
		apiURL = "https://rickandmortyapi.com/api/character"
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error al obtener personajes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResp struct {
		Results []map[string]interface{} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		http.Error(w, "Error al decodificar respuesta", http.StatusInternalServerError)
		return
	}

	var characters []Character
	for _, item := range apiResp.Results {
		char := Character{
			ID:      int(item["id"].(float64)),
			Name:    item["name"].(string),
			Status:  item["status"].(string),
			Image:   item["image"].(string),
			Gender:  item["gender"].(string),
			Species: item["species"].(string),
		}
		characters = append(characters, char)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
