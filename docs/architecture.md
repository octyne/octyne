# Architecture

## Overview

Octyne is an AI gateway with compatibility APIs at the edge, provider-neutral canonical models inside the system, and provider-specific adapters at the boundary to upstream AI services.

Current request flow:

```text
Client
-> Octyne HTTP server
-> compatibility handler
-> gateway service
-> model registry
-> provider registry
-> provider adapter
-> external provider HTTP API
-> response translation
-> client
```

The system is OpenAI-compatible first, but it must not become permanently OpenAI-centric.

## Package Responsibilities

- `cmd/octyne`: executable entry point. Load config, construct the app, start the server, and handle fatal startup errors.
- `internal/app`: composition root. Construct the dependency graph: provider configs, adapters, registries, gateway, and server.
- `internal/config`: load and validate environment configuration once during startup.
- `internal/server`: own routes, HTTP decoding, basic validation, compatibility response formatting, and status codes.
- `internal/gateway`: orchestrate chat requests, resolve models and providers, and delegate to adapters.
- `internal/registry`: map public model names to provider metadata.
- `internal/providers`: define configured upstream providers and the provider registry.
- `internal/adapters`: define adapter contracts.
- `internal/adapters/<provider>`: translate canonical requests to provider requests, call provider HTTP APIs, and translate responses back.
- `internal/types`: canonical provider-neutral request, response, and stream DTOs.

## Composition and Dependency Injection

Dependencies should be constructed explicitly in `internal/app`. Avoid package-level mutable globals and hidden initialization.

Preferred shape:

```go
providerRegistry := providers.NewRegistry()
gatewayService := gateway.New(providerRegistry)
httpServer := server.New(gatewayService)
```

Main should remain small. The server should not construct the gateway, the gateway should not construct registries, and adapters should not read environment variables.

## Schema Layers

Octyne must distinguish three schema layers:

```text
external compatibility schema
-> canonical Octyne model
-> provider-specific schema
```

On the response path:

```text
provider-specific schema
-> canonical Octyne model
-> requested compatibility schema
```

The current OpenAI-compatible request and canonical DTOs may look similar, but they should not be treated as permanently identical. New fields should be added only after deciding their canonical meaning, provider mappings, unsupported-provider behavior, and compatibility requirements.

The Chat Completions request boundary currently covers all 37 top-level OpenAI parameters. Stable nested shapes use typed structs and tagged unions; `json.RawMessage` is reserved for genuinely arbitrary user-provided JSON Schema and function-parameter payloads. Role-specific messages are decoded at the compatibility boundary and normalized into canonical chat messages before provider translation.

## Providers and Adapters

A provider is a configured upstream. An adapter is a protocol implementation. Multiple providers may share one adapter when their protocol is compatible.

Examples:

```text
OpenAI provider -> OpenAI-compatible adapter
Azure OpenAI provider -> OpenAI-compatible adapter
Ollama provider -> OpenAI-compatible adapter
```

Adapters should use provider-specific structs and direct HTTP calls. Prefer `net/http`, `encoding/json`, and `context`; avoid official or third-party provider SDKs unless a strong technical need appears.

Adapters should own reusable HTTP clients or receive injected clients for testing and future transport customization. Do not create a new `http.Client` per request.

## Context Propagation

Incoming request context must flow through the request path:

```text
r.Context()
-> server
-> gateway
-> adapter
-> outbound HTTP request
```

Use `http.NewRequestWithContext` for provider calls. Do not replace request context with `context.Background()` inside request handling. This preserves cancellation, timeouts, streaming shutdown, tracing, and resource cleanup.

## Credentials and BYOK Direction

Octyne credentials and provider credentials are separate. Provider API keys must not be placed in chat JSON bodies. Early configuration may use startup provider credentials, but future BYOK requires request- or account-scoped credential resolution.

Adapters should receive resolved credentials or typed auth options. They should not know whether credentials came from environment variables, a database, request metadata, a secret manager, or organization configuration. Do not log secrets.

## Streaming Design Notes

The current adapter interface includes `StreamChat`. Streaming should normalize provider events into canonical stream chunks, then the server should format those chunks for the requested compatibility API.

OpenAI-compatible streaming uses SSE:

```text
data: <json>

data: [DONE]
```

Streaming implementations must read incrementally, respect cancellation, avoid loading full streams into memory, close response bodies, close output channels exactly once, distinguish setup errors from in-stream errors, and flush each event to the client.

## Operational Direction

Move toward explicit `http.Server` timeouts, graceful shutdown, structured logging with `log/slog`, request IDs, bounded provider error body reads, and compatibility-layer error responses. Avoid exposing arbitrary internal Go error strings as the public API contract.
