package server

import (
	"encoding/json"
	"net/http"

	"github.com/usekeel/keel/internal/types"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req types.ChatCompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		http.Error(w, "model is required", http.StatusBadRequest)
		return
	}

	if len(req.Messages) == 0 {
		http.Error(w, "messages are required", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "chat completion not implemented",
	})
}
