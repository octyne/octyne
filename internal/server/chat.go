package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/types"
)

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	var compatReq openaicompat.ChatCompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&compatReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if compatReq.Model == "" {
		http.Error(w, "model is required", http.StatusBadRequest)
		return
	}

	if len(compatReq.Messages) == 0 {
		http.Error(w, "messages are required", http.StatusBadRequest)
		return
	}

	req := toCanonicalChatRequest(compatReq)

	if req.Stream {
		s.streamChat(w, r, req)
		return
	}

	resp, err := s.gateway.Chat(
		r.Context(),
		req,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) streamChat(
	w http.ResponseWriter,
	r *http.Request,
	req types.ChatCompletionRequest,
) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(
			w,
			"streaming is not supported",
			http.StatusInternalServerError,
		)
		return
	}

	chunks, err := s.gateway.StreamChat(
		r.Context(),
		req,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	for chunk := range chunks {
		if chunk.Error != nil {
			return
		}

		data, err := json.Marshal(chunk)
		if err != nil {
			return
		}

		if _, err := fmt.Fprintf(
			w,
			"data: %s\n\n",
			data,
		); err != nil {
			return
		}

		flusher.Flush()
	}

	if r.Context().Err() != nil {
		return
	}

	if _, err := fmt.Fprint(
		w,
		"data: [DONE]\n\n",
	); err != nil {
		return
	}

	flusher.Flush()
}
