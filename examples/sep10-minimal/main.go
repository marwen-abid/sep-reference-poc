package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"transaction": "example-challenge",
			"network_passphrase": "Test SDF Network ; September 2015",
		})
	})
	_ = http.ListenAndServe(":8081", nil)
}
