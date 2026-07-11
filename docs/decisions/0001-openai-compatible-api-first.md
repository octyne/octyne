# ADR 0001: OpenAI-Compatible API First

## Status

Accepted

## Context

Octyne aims to become a unified AI infrastructure platform, not only an LLM proxy. The first useful product slice needs to be easy for existing applications to adopt. Many clients and SDKs already support OpenAI-compatible chat completions, and many upstream services either expose OpenAI-compatible APIs or can be translated into that shape.

At the same time, Octyne must eventually support Anthropic-compatible and Gemini-compatible APIs, plus provider-neutral internal capabilities such as routing, policy, usage tracking, evals, and governance.

## Decision

Octyne will expose an OpenAI-compatible API first, starting with:

```http
POST /v1/chat/completions
```

Internally, Octyne will still preserve separate layers:

```text
external OpenAI-compatible schema
-> canonical Octyne model
-> provider-specific schema
```

OpenAI-compatible upstreams should reuse the same protocol implementation where practical. Anthropic and Gemini should first be exposed through Octyne's OpenAI-compatible API, with their own compatibility APIs added later.

## Consequences

Existing applications can adopt Octyne by changing base URL, API key, and model. Early development can produce value quickly with a small vertical slice.

The main risk is accidentally making the internal model permanently OpenAI-centric. Contributors must avoid copying every OpenAI field into canonical types without deciding whether it is provider-neutral, compatibility metadata, or provider-specific extension data.

## Implementation Notes

- Keep OpenAI HTTP schemas in `internal/adapters/openai`.
- Keep canonical DTOs in `internal/types`.
- Keep provider-specific translation inside adapters.
- Keep compatibility response formatting in `internal/server` until a dedicated compatibility layer is needed.
- Add fields incrementally with tests for mapping and unsupported-provider behavior.
