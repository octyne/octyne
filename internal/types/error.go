package types

import "fmt"

type ErrorKind string

const (
	ErrorKindInvalidRequest ErrorKind = "invalid_request"
	ErrorKindAuthentication ErrorKind = "authentication"
	ErrorKindPermission     ErrorKind = "permission"
	ErrorKindNotFound       ErrorKind = "not_found"
	ErrorKindRateLimit      ErrorKind = "rate_limit"
	ErrorKindTimeout        ErrorKind = "timeout"
	ErrorKindUnavailable    ErrorKind = "unavailable"
	ErrorKindInternal       ErrorKind = "internal"
)

// APIError is the provider-neutral error passed through Octyne's request path.
// ProviderRequestID is retained for diagnostics and must not replace Octyne's
// own client-facing request ID.
type APIError struct {
	Kind              ErrorKind
	Message           string
	Param             *string
	Code              *string
	HTTPStatus        int
	ProviderRequestID string
	Cause             error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Cause
}
