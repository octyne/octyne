# Octyne Project Context

## Identity

Octyne is a production-grade AI infrastructure platform written in Go. It was previously named Keel. The first product phase is a lightweight, high-performance AI gateway and proxy that exposes an OpenAI-compatible API so clients can change only `base_url`, API key, and model to route through Octyne.

The long-term goal is broader than an LLM proxy. Octyne should grow into a unified AI infrastructure platform for provider routing, model management, API keys, BYOK credentials, usage tracking, cost accounting, rate limiting, failover, load balancing, tracing, observability, prompt management, evaluations, guardrails, governance, and multiple API compatibility layers.

## Development Philosophy

Octyne should be developed incrementally, but every architectural decision should be suitable for a real product. Prefer small vertical or architectural increments:

```text
new branch -> one small change -> compile -> test -> manually verify -> commit -> pull request to main
```

Code changes should happen on focused branches and be merged to `main` through pull requests. Use idiomatic Go, clear package boundaries, explicit dependency injection, context propagation, testability, and maintainable interfaces. Avoid tutorial shortcuts, broad speculative refactors, unnecessary abstractions, and provider SDK dependencies.

The primary developer is experienced in Python and AI engineering and is learning Go through this project. When introducing Go-specific concepts or design choices, briefly explain why they fit.

## Product Direction

Phase 1 exposes an OpenAI-compatible API, starting with:

```http
POST /v1/chat/completions
```

OpenAI, Anthropic, and Gemini should initially be reachable through the OpenAI-compatible Octyne API. OpenAI-compatible upstreams such as Azure OpenAI, OpenRouter, Ollama, vLLM, and LM Studio should reuse the same protocol implementation where practical rather than copying adapter code.

Long term, Octyne should expose multiple compatibility APIs, including OpenAI-compatible, Anthropic-compatible, and Gemini-compatible APIs. External compatibility formats must remain separate from internal canonical models.

## Current State

The current working vertical slice supports both non-streaming and streaming chat completions through the OpenAI adapter:

```text
Client -> HTTP server -> chat handler -> gateway -> model registry
-> provider registry -> OpenAI adapter -> OpenAI API
-> translated canonical response -> client
```

The model registry is constructed in `internal/app` and injected into the
gateway. Each public model name resolves to a provider and an upstream model ID.
The gateway forwards that upstream ID for both non-streaming and streaming
requests. Public model names use the required `provider/model` format so clients
select a provider explicitly even when multiple providers offer the same model.

Current routes:

```http
GET /health
GET /v1/models
POST /v1/chat/completions
```

`GET /v1/models` reads the application-owned registry directly and returns the
OpenAI-compatible model-list envelope. Model IDs are the public
`provider/model` names clients use for routing, and `owned_by` identifies the
configured provider. Entries are sorted by public model ID for deterministic
responses.

Current low-cost OpenAI development model:

```text
openai/gpt-5-nano
```

For `POST /v1/chat/completions`, requests with `stream` omitted or set to `false` return OpenAI-compatible `chat.completion` JSON. Requests with `stream: true` return OpenAI-compatible SSE events in the form `data: <chat.completion.chunk JSON>\n\n`, followed by `data: [DONE]\n\n` after successful completion.

The OpenAI adapter shares request construction and provider routing across both modes while keeping response handling separate. Non-streaming decodes one JSON response. Streaming reads SSE incrementally, propagates request cancellation, and gives the stream goroutine sole ownership of closing the upstream body and output channel.

The current OpenAI non-streaming request timeout is 600 seconds. Streaming has no total-duration timeout; it uses a 30-second response-header timeout and remains governed by the incoming request context after the stream begins.

Automated tests cover typed non-streaming response translation, assistant outputs, token accounting, log probabilities, moderation and service metadata, streaming deltas, multiple choices, explicit-null and usage-only chunks, downstream SSE framing, upstream stream parsing, `[DONE]`, malformed chunks, provider setup errors, cancellation, timeout behavior, and response-body closure. Default tests use local HTTP test servers and do not call paid provider APIs.

The OpenAI-compatible Chat Completions request now covers all 37 current top-level parameters, including typed role-specific and multimodal messages, tools and tool choices, structured output, prediction, streaming options, prompt caching, provider-assisted features, and accepted deprecated fields. The compatibility, canonical, and OpenAI provider layers remain distinct, and optional scalar values preserve explicit zero, false, and empty values.

The current documented non-streaming response and streaming chunk schemas are typed through the OpenAI provider and canonical layers. Chat Completions errors use OpenAI-compatible envelopes and status mapping, every response receives an Octyne request ID, and the same ID is forwarded to OpenAI for correlation without replacing the provider's own diagnostic request ID. Upstream server details are sanitized, and failures after streaming headers are committed are returned as SSE error events without a successful `[DONE]` terminator.

## Near-Term Priorities

1. Move startup model registrations from the composition root toward configuration-driven registration.
2. Add explicit server timeouts, graceful shutdown, and structured logging.
3. Add Octyne API authentication separate from provider credentials.
4. Begin reusable OpenAI-compatible provider configuration.

## Beta Scope

Beta should focus on a dependable unified AI gateway:

- OpenAI, Anthropic, and Gemini support through the OpenAI-compatible API.
- Non-streaming and streaming chat.
- Model registry and `GET /v1/models`.
- Octyne API authentication separate from provider credentials.
- BYOK and Octyne-managed provider credentials.
- Consistent compatibility-layer error responses.
- Structured logging, request IDs, timeouts, and graceful shutdown.
- Tests that avoid paid provider APIs by default.

## Long-Term Shape

Octyne should eventually separate into a lightweight data plane and a control plane.

The data plane handles latency-sensitive request flow: auth, validation, routing, provider calls, streaming, rate limiting, usage collection, failover, caching, and policy enforcement.

The control plane handles management: users, organizations, projects, API keys, provider credentials, models, routing rules, prompts, evals, traces, billing, policies, and audit logs.

## Repository Memory Rules

Keep durable project memory in `docs/project-context.md`, stable technical design in `docs/architecture.md`, milestone sequencing in `docs/roadmap.md`, and major decisions in `docs/decisions/`. Avoid creating a loose memory folder unless the repository later needs machine-readable state.
