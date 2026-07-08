# Octyne

> Unified AI Gateway for every model provider.

Octyne is an OpenAI-compatible AI gateway that provides a single API for multiple AI providers.

Instead of integrating separately with OpenAI, Anthropic, Google, Groq, or local models, applications integrate with Octyne once. Octyne handles provider routing, model abstraction, authentication, observability, and more.

The project is designed to evolve into a complete AI infrastructure platform with routing, governance, analytics, caching, guardrails, and model management.

---

## Features

Current

- OpenAI-compatible Chat Completions API
- Provider abstraction layer
- Model registry
- Provider registry
- HTTP-based providers (no provider SDKs)
- Written in Go
- Lightweight and dependency minimal

Planned

- Streaming responses
- Embeddings API
- Image generation
- Audio APIs
- Responses API
- Provider failover
- Smart routing
- Load balancing
- BYOK (Bring Your Own Keys)
- Multi-tenancy
- API Keys
- Rate limiting
- Usage analytics
- OpenTelemetry
- Prompt caching
- Guardrails
- Model marketplace

---

## Why Octyne?

Every AI provider exposes slightly different APIs.

Octyne hides those differences behind one consistent interface.

```
                Your Application
                       │
                       ▼
                 ┌────────────┐
                 │  Octyne    │
                 └────────────┘
          ┌────────┼────────┬────────┐
          ▼        ▼        ▼        ▼
       OpenAI   Anthropic Gemini  Ollama
```

Switch providers without changing application code.

---

## Project Structure

```
cmd/
    octyne/

internal/
    adapters/
    app/
    config/
    gateway/
    providers/
    registry/
    server/
    types/
```

---

## Requirements

- Go 1.26+
- OpenAI API Key

---

## Installation

Clone the repository.

```bash
git clone https://github.com/<username>/octyne.git
cd octyne
```

Install dependencies.

```bash
go mod tidy
```

---

## Configuration

Create a `.env` file.

```env
OPENAI_API_KEY=your_api_key
PORT=3000
```

---

## Running

Start the server.

```bash
go run ./cmd/octyne
```

You should see:

```
2026/07/08 12:36:55 Octyne starting on :3000
```

---

## Example Request

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model":"gpt-5-nano",
    "messages":[
      {
        "role":"user",
        "content":"hello"
      }
    ]
  }'
```

## Example response

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

---

# API

## Chat Completions

### Endpoint

```
POST /v1/chat/completions
```

### Request

```json
{
  "model": "gpt-5-nano",
  "messages": [
    {
      "role": "user",
      "content": "Write a haiku"
    }
  ]
}
```

---

# Supported Models

Currently

- all openai models compatible with chat completions API are supported.

More providers and models will be added over time.

---

# Development

Run

```bash
go run ./cmd/octyne
```

Build

```bash
go build ./cmd/octyne
```

---

# Design Principles

- OpenAI compatible
- Provider agnostic
- Cloud agnostic
- Minimal dependencies
- No provider SDKs
- Fast startup
- Production-ready architecture
- Extensible provider system

---

# Roadmap

- [x] OpenAI Chat Completions
- [ ] Streaming
- [ ] Embeddings
- [ ] Images
- [ ] Audio
- [ ] Anthropic provider
- [ ] Gemini provider
- [ ] Ollama provider
- [ ] OpenRouter provider
- [ ] API Keys
- [ ] Authentication
- [ ] Multi-tenancy
- [ ] Metrics
- [ ] OpenTelemetry
- [ ] Rate limiting
- [ ] Model routing
- [ ] Retries
- [ ] Fallback providers
- [ ] Prompt caching
- [ ] Dashboard
- [ ] SDKs

---

# FAQ

### Why not call providers directly?

Every provider has different APIs, authentication methods, models, and response formats. Octyne provides a single interface so your application doesn't need provider-specific integrations.

---

### Is Octyne an OpenAI replacement?

No.

Octyne sits between your application and AI providers. It forwards requests to supported providers while exposing a consistent API.

---

### Does Octyne use provider SDKs?

No.

Octyne communicates with providers over standard HTTP APIs to remain lightweight and provider-independent.

---

### Can I add my own provider?

Yes.

Octyne is built around an adapter architecture, making it straightforward to add new providers without changing the gateway.

---

### Is Octyne production ready?

The architecture is designed for production use, but the project is currently in active development. New features and providers are being added incrementally.

---

# Contributing

Contributions, issues, and discussions are welcome.

---

# License

MIT