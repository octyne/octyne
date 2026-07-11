# Octyne

Octyne is a lightweight AI gateway written in Go.

Today it exposes an OpenAI-compatible Chat Completions endpoint and routes requests through the configured OpenAI provider. The project is designed to grow into a broader AI infrastructure platform for routing, credentials, usage tracking, observability, governance, and multiple provider APIs.

## Current Support

- `POST /v1/chat/completions`
- `GET /health`
- OpenAI provider adapter
- In-memory model registry
- Provider abstraction layer
- Docker and Compose local runtime

Streaming and additional providers are planned, but not enabled yet.

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

## Chat Example

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5-nano",
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
  "model": "gpt-5-nano",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Hello!"
      }
    }
  ]
}
```

## Supported Models

Currently registered models:

- `gpt-5-nano`
- `gpt-4.1-mini`

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
- `internal/registry`: model-to-provider mappings
- `internal/types`: provider-neutral DTOs

## Project Docs

- [Project context](docs/project-context.md)
- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)
- [Decisions](docs/decisions/)
