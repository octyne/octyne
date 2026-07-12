package server

import (
	"encoding/json"
	"net/http"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
)

func (s *Server) modelsHandler(w http.ResponseWriter, r *http.Request) {
	registeredModels := s.modelRegistry.List()
	models := make([]openaicompat.Model, len(registeredModels))

	for i, registeredModel := range registeredModels {
		models[i] = openaicompat.Model{
			ID:      registeredModel.Name,
			Object:  "model",
			Created: 0,
			OwnedBy: registeredModel.Model.Provider,
		}
	}

	response := openaicompat.ModelList{
		Object: "list",
		Data:   models,
	}

	data, err := json.Marshal(response)
	if err != nil {
		writeOpenAIError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(append(data, '\n'))
}
