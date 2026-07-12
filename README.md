# Octyne

Octyne is a lightweight AI gateway written in Go.

Today it exposes an OpenAI-compatible Chat Completions endpoint and routes requests through the configured OpenAI provider. The project is designed to grow into a broader AI infrastructure platform for routing, credentials, usage tracking, observability, governance, and multiple provider APIs.

## Current Support

- `POST /v1/chat/completions`
- `GET /health`
- Non-streaming OpenAI-compatible chat completions
- Streaming OpenAI-compatible chat completions over SSE
- Complete typed Chat Completions requests, success responses, and streaming chunks
- OpenAI-compatible error envelopes and per-request `x-request-id` headers
- OpenAI provider adapter
- Injected in-memory model registry with public-to-upstream model mapping
- Provider abstraction layer
- Docker and Compose local runtime

Additional providers are planned, but not enabled yet.

## Requirements

- Go 1.26+
- OpenAI API key

## Configuration

Create a `.env` file:

```env
OPENAI_API_KEY=your_api_key
PORT=3000
```

`PORT` is optional and defaults to `3000`.

## Run Locally

```bash
go run ./cmd/octyne
```

Or use the Makefile:

```bash
make run
```

Health check:

```bash
curl http://localhost:3000/health
```

## Chat Examples

### Non-streaming

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai/gpt-5-nano",
    "messages": [
      {
        "role": "user",
        "content": "hello"
      }
    ]
  }'
```

Example response shape:

```json
{
  "id": "...",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "gpt-5-nano",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello!",
        "refusal": null
      },
      "finish_reason": "stop",
      "logprobs": null
    }
  ],
  "usage": {
    "prompt_tokens": 8,
    "completion_tokens": 2,
    "total_tokens": 10
  }
}
```

The response `model` value is the upstream model identifier reported by the
selected provider.

### Streaming

Set `stream` to `true` to receive OpenAI-compatible server-sent events. A
successful stream ends with `data: [DONE]`.

```bash
curl -N http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai/gpt-5-nano",
    "messages": [
      {
        "role": "user",
        "content": "hello"
      }
    ],
    "stream": true
  }'
```

## Request Parameter Compatibility

The OpenAI-compatible request boundary supports the current Chat Completions
request schema, including multimodal messages, tools, structured output, and
accepted deprecated fields. Optional scalar parameters preserve explicit zero,
`false`, and empty values while requests move through the compatibility,
canonical, and provider layers. See the
[Chat Completions schema status](docs/chat-completions-schema.md) for the full
parameter list and nested-shape coverage.

Deprecated OpenAI parameters remain accepted for compatibility:

- `max_tokens` (prefer `max_completion_tokens`)
- `user` (prefer `safety_identifier` and/or `prompt_cache_key`, depending on the
  use case)
- `prompt_cache_retention` (prefer `prompt_cache_options.ttl`)
- `functions` (prefer `tools`)
- `function_call` (prefer `tool_choice`)

For the OpenAI provider, Octyne currently preserves and forwards these fields.
New clients should use the replacement fields. Future provider adapters will
translate a deprecated field only when that provider has a safe equivalent;
otherwise Octyne should return a clear compatibility error rather than silently
changing or dropping the request. Provider-native compatibility APIs should
expose that API's own parameter names and deprecation rules instead of inheriting
OpenAI-only legacy fields.

## Response Compatibility

Successful non-streaming responses preserve the current typed Chat Completions
schema, including nullable assistant content, refusals, URL citations, audio,
function and custom tool calls, deprecated function calls, finish reasons,
content and refusal log probabilities, complete token usage details, moderation,
service tier, and system fingerprint metadata.

Streaming responses preserve typed content, refusal, deprecated function-call,
and indexed function tool-call fragments; multiple choices; log probabilities;
moderation and service metadata; obfuscation; explicit-null intermediate usage;
and the final usage-only chunk. A successful stream still ends with
`data: [DONE]`. See the
[Chat Completions schema status](docs/chat-completions-schema.md) for the full
response inventory.

## Error Compatibility and Request IDs

Chat Completions errors use the OpenAI-compatible JSON envelope with `message`,
`type`, nullable `param`, and nullable `code` fields. Validation, routing,
provider status, rate-limit, timeout, and internal failures are mapped to safe
HTTP statuses and public messages instead of exposing arbitrary Go or upstream
server details.

Every Chat Completions response includes an Octyne-generated `x-request-id`.
Octyne also sends that value to the OpenAI provider as `X-Client-Request-Id` for
cross-system correlation while retaining the provider's own request ID as
internal diagnostic metadata. If a stream fails after its HTTP headers have
already been sent, Octyne emits an error envelope as an SSE `data:` event and
does not emit `[DONE]`.

## Supported Models

Currently registered models:

- `openai/gpt-5-nano`
- `openai/gpt-4.1-mini`

Public model names use the required `provider/model` format. Registry entries
resolve those names to a provider and its upstream model ID, so clients can
select the provider explicitly while adapters receive provider-native IDs.

## Docker

```bash
docker compose up --build
```

## Development

```bash
make test
make vet
make check
```

Useful package boundaries:

- `internal/server`: HTTP routes and compatibility response formatting
- `internal/gateway`: request orchestration and provider resolution
- `internal/adapters/openai`: OpenAI request/response translation
- `internal/providers`: configured upstream providers
- `internal/registry`: public-model-to-provider and upstream-model mappings
- `internal/types`: provider-neutral DTOs

## Project Docs

- [Project context](docs/project-context.md)
- [Architecture](docs/architecture.md)
- [Chat Completions schema status](docs/chat-completions-schema.md)
- [Roadmap](docs/roadmap.md)
- [Decisions](docs/decisions/)
