# Repository Guidelines

## Project Identity & Product Direction

Octyne is an AI infrastructure platform written in Go. The first phase is a lightweight, high-performance AI gateway exposing an OpenAI-compatible API, but the long-term product is broader: provider routing, model management, API keys, BYOK credentials, usage and cost tracking, rate limiting, failover, observability, prompt management, evals, guardrails, governance, and multiple compatibility APIs.

## Project Memory & Architecture Docs

Use `docs/project-context.md` for durable product context, current state, and long-term memory. Use `docs/architecture.md` for system design, package responsibilities, schema boundaries, provider rules, credentials, and streaming notes. Use `docs/roadmap.md` for milestone sequencing. Record major architectural choices in `docs/decisions/` as ADRs. Read these docs before significant feature work or refactors, and update them when decisions or direction change.

## Project Structure & Module Organization

The executable entry point is `cmd/octyne/main.go`; keep it small. `internal/app` is the composition root and should construct config, providers, adapters, gateway, and server dependencies. `internal/server` owns HTTP routes and compatibility response formatting. `internal/gateway` orchestrates model/provider resolution and delegates to adapters. `internal/providers` defines configured upstream providers. `internal/adapters/<provider>` translates and calls provider HTTP APIs. `internal/registry` maps public models to providers. `internal/types` holds canonical provider-neutral DTOs.

## Build, Test, and Development Commands

- `go run ./cmd/octyne`: start the local server. Requires `.env` values such as `OPENAI_API_KEY` and optionally `PORT`.
- `go build ./cmd/octyne`: compile the Octyne binary.
- `go test ./...`: run all package tests; use this before opening a PR even when adding only one package.
- `go vet ./...`: catch common Go correctness issues before commits.
- `go mod tidy`: reconcile `go.mod` and `go.sum` after dependency changes.
- `gofmt -w <files>`: format changed Go files before committing.

## Architecture Principles

Use idiomatic Go and explain new Go concepts when they affect design. Prefer small, explicit interfaces, dependency injection, package boundaries, and composition. Avoid package-level mutable globals. Propagate `r.Context()` through server, gateway, adapter, and outbound HTTP calls; do not replace request context with `context.Background()` in request paths. Use values for DTOs and pointers for stateful services.

Keep three schema layers distinct: external compatibility schemas, canonical Octyne types, and provider-specific schemas. The initial OpenAI-compatible API must not make the internal model permanently OpenAI-centric.

## Provider, Credential & HTTP Rules

Prefer the Go standard library (`net/http`, `encoding/json`, `context`) and avoid provider SDKs unless there is a strong technical reason. Reuse HTTP clients; do not create one per request. Adapters should receive configuration or resolved credentials and must not call `os.Getenv()` directly. Never put provider API keys in chat JSON bodies, never log secrets, and keep future BYOK support in mind when designing credential flow.

## Coding Style & Naming Conventions

Follow `gofmt`, short package names, exported identifiers only for cross-package APIs, and wrapped errors with operation context. Keep provider-specific translation in `internal/adapters/<provider>` and shared contracts in `internal/types`, `internal/providers`, or focused gateway interfaces. Add abstractions only when they represent confirmed long-term concepts or remove real duplication.

## Testing Guidelines

Add package-local `*_test.go` files. Prefer table-driven tests for translation, config validation, registries, routing, and error paths. Use fake adapters for gateway tests and `httptest.Server` for provider adapters. Default tests must not call paid provider APIs; live provider checks should be opt-in. Use `gpt-5-nano` for routine low-cost OpenAI smoke tests.

## Commit & Pull Request Guidelines

Do code changes on focused branches and raise pull requests into `main`; do not commit code directly to `main`. Use Conventional Commits for all commit messages, such as `docs: update README`, `feat: add streaming adapter contract`, or `refactor: rename project from Keel to Octyne`. Keep commits focused, brief, and deployable; each commit should be safe to ship independently. When suggesting commits, provide a Conventional Commit subject plus a brief body that explains what changed and why. Do not add Codex, AI, or tool-generated signatures to commit messages, descriptions, or PR text. PRs should include a summary, test results, linked issues when applicable, and example requests or logs for API behavior changes.

## Current Development State

The current vertical slice is non-streaming `POST /v1/chat/completions` through the OpenAI adapter. The next priority is OpenAI streaming: extract shared OpenAI chat request construction, force `stream: true` in `StreamChat`, parse SSE incrementally, propagate cancellation, close channels and bodies exactly once, add gateway/server streaming paths, and return OpenAI-compatible SSE from the same endpoint when `req.Stream` is true. Preserve non-streaming behavior while implementing streaming.

## Security & Configuration Tips

Do not commit real `.env` values or API keys. Document any new environment variable in `README.md`, and keep provider credentials flowing through configuration rather than hard-coded defaults.
