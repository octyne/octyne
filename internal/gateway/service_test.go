package gateway

import (
	"errors"
	"net/http"
	"testing"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

func TestResolveAdapterReturnsTypedRoutingErrors(t *testing.T) {
	tests := []struct {
		name       string
		model      string
		registry   *providers.Registry
		wantKind   types.ErrorKind
		wantStatus int
		wantParam  string
	}{
		{
			name:       "unknown model",
			model:      "not-registered",
			registry:   providers.NewRegistry(),
			wantKind:   types.ErrorKindNotFound,
			wantStatus: http.StatusNotFound,
			wantParam:  "model",
		},
		{
			name:       "provider missing",
			model:      "gpt-5-nano",
			registry:   providers.NewRegistry(),
			wantKind:   types.ErrorKindInternal,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "adapter missing",
			model: "gpt-5-nano",
			registry: registryWithProvider(providers.New(
				providers.Config{Name: "openai"},
				nil,
			)),
			wantKind:   types.ErrorKindInternal,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.registry).resolveAdapter(tt.model)
			var apiErr *types.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %T, want *types.APIError", err)
			}
			if apiErr.Kind != tt.wantKind || apiErr.HTTPStatus != tt.wantStatus {
				t.Errorf("error = %+v, want kind %q status %d", apiErr, tt.wantKind, tt.wantStatus)
			}
			if stringValue(apiErr.Param) != tt.wantParam {
				t.Errorf("Param = %q, want %q", stringValue(apiErr.Param), tt.wantParam)
			}
		})
	}
}

func registryWithProvider(provider *providers.Provider) *providers.Registry {
	registry := providers.NewRegistry()
	registry.Register(provider.Name, provider)
	return registry
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
