# Roadmap

This roadmap is directional. Keep implementation incremental and commit each coherent change separately.

## Milestone 1: Complete OpenAI Chat

The initial non-streaming and streaming vertical slice is complete.

- [x] Preserve the working non-streaming chat completions path.
- [x] Extract shared OpenAI chat request construction.
- [x] Implement OpenAI streaming request execution.
- [x] Parse OpenAI SSE incrementally, including `[DONE]`.
- [x] Add gateway and server streaming paths.
- [x] Return OpenAI-compatible SSE when `stream: true`.
- [x] Add focused tests for non-streaming behavior, streaming parsing, provider setup errors, cancellation, timeout behavior, channel closure, and downstream SSE framing.
- [x] Complete all current top-level Chat Completions request parameters and their typed nested request shapes.
- [x] Complete typed non-streaming response messages, finish reasons, log probabilities, usage, moderation, and service metadata.
- [x] Complete typed streaming deltas, tool-call fragments, multiple choices, metadata, obfuscation, and usage-only chunks.
- [ ] Add request IDs and canonical OpenAI-compatible errors.
- [ ] Expand focused tests for remaining translation, routing, and configuration paths.

## Milestone 2: OpenAI-Compatible Providers

Reuse the OpenAI-compatible adapter where practical for:

- OpenAI
- Azure OpenAI
- OpenRouter
- Ollama
- vLLM
- LM Studio

Provider-specific differences should be configuration or focused extensions, not copied adapters. Differences may include base URL, authentication, headers, query parameters, deployment paths, API versions, and unsupported fields.

## Milestone 3: Anthropic Adapter

Expose Anthropic models through the OpenAI-compatible API first.

- Request translation.
- System message handling.
- Chat message conversion.
- Tool-call mapping.
- Streaming event translation.
- Usage mapping.
- Finish reason mapping.
- Provider error mapping.

## Milestone 4: Gemini Adapter

Expose Gemini models through the OpenAI-compatible API first.

- Role and content translation.
- System instruction mapping.
- Generation config mapping.
- Tool mapping.
- Streaming.
- Usage mapping.
- Finish reason mapping.
- Provider error mapping.

## Milestone 5: Authentication and BYOK

- Octyne API keys.
- Hashed key storage.
- Request authentication middleware.
- Provider credential selection.
- BYOK header or stored credential flow.
- Secret encryption.
- Credential resolver design.

Octyne API keys authenticate clients to Octyne. Provider credentials authenticate Octyne to upstream providers. Keep these concerns separate.

## Milestone 6: Model API and Registry Evolution

- Add `GET /v1/models`.
- Return OpenAI-compatible model listings.
- Move from hardcoded in-memory models toward configuration-driven registration.
- Preserve room for future metadata: provider mapping, aliases, capabilities, pricing, context window, streaming support, tools, vision, availability, routing policy, deployment ID, and organization visibility.

## Milestone 7: Operational Readiness

- Structured logging.
- Metrics.
- Request IDs.
- Explicit server timeouts.
- Graceful shutdown.
- Retry policy where safe.
- Circuit breaking where appropriate.
- Rate limiting.
- Docker image.
- CI.
- Release builds.
- Configuration documentation.

Retries require care for non-idempotent and streaming operations.

## Later Platform Capabilities

- Routing policies: explicit model, aliases, cheapest, lowest latency, provider priority, weighted routing, geographic routing, availability-based routing, capability-based routing, organization rules, and fallback chains.
- Observability: traces, provider latency, time to first token, total latency, token usage, error rates, routing decisions, cost, retry count, and fallback count.
- Prompt management: versions, environments, variables, deployment history, rollback, aliases, and experiments.
- Evaluations: datasets, offline and online evals, regression testing, model comparison, and trace-linked evaluation.
- Usage and billing: token accounting, cost accounting, quotas, budgets, alerts, and asynchronous usage events.
- Multi-tenancy: users, organizations, projects, API keys, provider credentials, policies, and isolation.
- Security: key hashing, secret encryption, rotation, audit logs, least privilege, header redaction, prompt/completion redaction, SSRF protection, and allowed-host policies.
