package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/types"
)

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	var compatReq openaicompat.ChatCompletionRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&compatReq); err != nil {
		writeInvalidJSONError(w, err)
		return
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		writeInvalidJSONError(w, err)
		return
	}

	if compatReq.Model == "" {
		param := "model"
		writeOpenAIError(w, &types.APIError{
			Kind:       types.ErrorKindInvalidRequest,
			Message:    "Missing required parameter: 'model'.",
			Param:      &param,
			HTTPStatus: http.StatusBadRequest,
		})
		return
	}

	if len(compatReq.Messages) == 0 {
		param := "messages"
		writeOpenAIError(w, &types.APIError{
			Kind:       types.ErrorKindInvalidRequest,
			Message:    "Missing required parameter: 'messages'.",
			Param:      &param,
			HTTPStatus: http.StatusBadRequest,
		})
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
		writeOpenAIError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		writeOpenAIError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(append(data, '\n'))
}

func writeInvalidJSONError(w http.ResponseWriter, cause error) {
	code := "invalid_json"
	writeOpenAIError(w, &types.APIError{
		Kind:       types.ErrorKindInvalidRequest,
		Message:    "Invalid request body.",
		Code:       &code,
		HTTPStatus: http.StatusBadRequest,
		Cause:      cause,
	})
}

func (s *Server) streamChat(
	w http.ResponseWriter,
	r *http.Request,
	req types.ChatCompletionRequest,
) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeOpenAIError(w, &types.APIError{
			Kind:       types.ErrorKindInternal,
			Message:    "Streaming is not supported.",
			HTTPStatus: http.StatusInternalServerError,
		})
		return
	}

	chunks, err := s.gateway.StreamChat(
		r.Context(),
		req,
	)
	if err != nil {
		writeOpenAIError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	for chunk := range chunks {
		if chunk.Error != nil {
			data, err := json.Marshal(
				toOpenAIErrorEnvelope(normalizeAPIError(chunk.Error)),
			)
			if err == nil {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
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
