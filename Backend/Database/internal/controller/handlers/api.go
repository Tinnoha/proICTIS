package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	config := map[string]string{
		"base_url": baseURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
