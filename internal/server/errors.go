package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/octyne/octyne/internal/types"
)

type openAIErrorEnvelope struct {
	Error openAIError `json:"error"`
}

type openAIError struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    *string `json:"code"`
}

func writeOpenAIError(w http.ResponseWriter, err error) {
	apiErr := normalizeAPIError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.HTTPStatus)
	_ = json.NewEncoder(w).Encode(toOpenAIErrorEnvelope(apiErr))
}

func normalizeAPIError(err error) *types.APIError {
	var apiErr *types.APIError
	if errors.As(err, &apiErr) {
		if apiErr.HTTPStatus >= 400 && apiErr.HTTPStatus <= 599 {
			return apiErr
		}
	}

	return &types.APIError{
		Kind:       types.ErrorKindInternal,
		Message:    "Internal server error.",
		HTTPStatus: http.StatusInternalServerError,
		Cause:      err,
	}
}

func toOpenAIErrorEnvelope(err *types.APIError) openAIErrorEnvelope {
	errorType := "server_error"
	switch err.Kind {
	case types.ErrorKindInvalidRequest, types.ErrorKindNotFound:
		errorType = "invalid_request_error"
	case types.ErrorKindAuthentication:
		errorType = "authentication_error"
	case types.ErrorKindPermission:
		errorType = "permission_error"
	case types.ErrorKindRateLimit:
		errorType = "rate_limit_error"
	}

	return openAIErrorEnvelope{Error: openAIError{
		Message: err.Message,
		Type:    errorType,
		Param:   err.Param,
		Code:    err.Code,
	}}
}
