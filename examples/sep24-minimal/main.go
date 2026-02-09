package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/sep24/info", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"deposit": map[string]any{"USDC": map[string]any{"enabled": true}},
			"withdraw": map[string]any{"USDC": map[string]any{"enabled": true}},
		})
	})
	_ = http.ListenAndServe(":8082", nil)
}
