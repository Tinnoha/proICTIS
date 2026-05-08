package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := map[string]string{
		"base_url": os.Getenv("BASE_URL"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
