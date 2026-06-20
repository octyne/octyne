package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func healthHandler(w http.ResponseWriter, r *http.Request){
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err!=nil{
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	log.Println("Keel starting on :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
