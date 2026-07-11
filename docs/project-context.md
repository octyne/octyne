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

The current working vertical slice is non-streaming chat completions through the OpenAI adapter:

```text
Client -> HTTP server -> chat handler -> gateway -> model registry
-> provider registry -> OpenAI adapter -> OpenAI API
-> translated canonical response -> client
```

Current routes:

```http
GET /health
POST /v1/chat/completions
```

Current low-cost OpenAI development model:

```text
gpt-5-nano
```

The next implementation priority is OpenAI streaming support while preserving existing non-streaming behavior.

## Near-Term Priorities

1. Complete OpenAI streaming for `POST /v1/chat/completions`.
2. Add common generation parameters carefully after streaming is stable.
3. Improve canonical error handling and OpenAI-compatible error responses.
4. Add focused tests for translation, routing, config, handlers, and adapters.
5. Move the model registry toward configurable registration before it becomes permanent hardcoding.

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
