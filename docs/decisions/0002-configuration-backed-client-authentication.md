# ADR 0002: Configuration-Backed Client Authentication

## Status

Accepted

## Context

Octyne's OpenAI-compatible endpoints need to authenticate clients independently
of the credentials Octyne uses with upstream providers. The first increment
needs to be small and secure without introducing users, organizations, a
database, key-management APIs, or BYOK credential resolution prematurely.

## Decision

Require one or more client keys through `OCTYNE_API_KEYS` at startup. Protect
the exact `/v1` path and the `/v1/*` route group with bearer authentication,
while keeping `GET /health` public.

Construct a static verifier in the application composition root. It retains
SHA-256 digests of configured keys and compares a presented key against every
digest with constant-time comparison. The HTTP middleware accepts exactly one
`Authorization: Bearer <key>` credential and returns an OpenAI-compatible
`401 Unauthorized` error for missing, malformed, or invalid credentials.

Keep Octyne client keys outside provider configuration. Provider credentials
continue to flow only from provider configuration into provider adapters, and
the incoming client `Authorization` header is never forwarded upstream.

## Consequences

- `/v1` and future `/v1/*` endpoints are protected by default.
- `/health` remains available to unauthenticated health probes.
- Multiple configured keys permit separate clients and rotation overlap.
- Request IDs and structured logs cover rejected requests without logging keys.
- Restarting Octyne is required to add or revoke a key.
- Plaintext keys still originate in process configuration; this is not
  persistent hashed key storage.
- Persistent key ownership, identifiers, revocation, rotation APIs, audit
  events, rate limiting, secret encryption, and BYOK remain future work.
