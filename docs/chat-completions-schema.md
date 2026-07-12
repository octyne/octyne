# Chat Completions Schema Status

Snapshot date: 2026-07-12.

The OpenAI-compatible `POST /v1/chat/completions` request implements the
top-level request fields documented in the current
[OpenAI Chat Completions reference](https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create)
through the compatibility, canonical, and OpenAI provider layers. No documented
top-level create parameter is known to be missing as of this snapshot.

This is a point-in-time implementation inventory, not a permanent parameter-count
guarantee. Recheck the upstream reference when OpenAI changes the API.

## Completed top-level request parameters

- model
- messages
- stream
- temperature
- top_p
- frequency_penalty
- presence_penalty
- max_completion_tokens
- n
- logprobs
- top_logprobs
- reasoning_effort
- verbosity
- seed
- store
- parallel_tool_calls
- safety_identifier
- prompt_cache_key
- max_tokens
- user
- prompt_cache_retention
- metadata
- service_tier
- prompt_cache_options
- stop
- logit_bias
- stream_options
- modalities
- audio
- response_format
- prediction
- moderation
- web_search_options
- tools
- tool_choice
- functions
- function_call

## Deprecated request parameters

Octyne continues to accept and forward these OpenAI fields for compatibility:

- `max_tokens`: use `max_completion_tokens` for new clients.
- `user`: use `safety_identifier` for abuse detection and `prompt_cache_key`
  for cache bucketing, as applicable.
- `prompt_cache_retention`: use `prompt_cache_options.ttl`.
- `functions`: use `tools`.
- `function_call`: use `tool_choice`.

Accepting a field at the compatibility boundary does not mean every model
supports it. OpenAI documents model-specific differences, especially for
reasoning models. Octyne preserves the request shape; the selected provider and
model remain authoritative for feature support.

The request schema includes role-specific developer, system, user, assistant, tool, and deprecated function messages; text, image, input-audio, file, and refusal content; assistant audio references; function and custom tool calls; tool results; typed response formats; typed streaming and prompt-cache options; and explicit assistant null content.

## Remaining response-schema typing

Streaming and non-streaming request execution are complete. The remaining work
is to replace the intentionally minimal response DTOs with complete typed
OpenAI-compatible response shapes:

- Complete non-streaming response messages, finish reasons, annotations, audio, tool calls, and deprecated function calls.
- Replace response-side `any` fields with typed usage and content/refusal log probability structures.
- Complete streaming deltas for refusals, audio, annotations, function/custom tool-call fragments, multiple choices, usage-only chunks, and obfuscation.
- Preserve successful `[DONE]` behavior and cancellation semantics while expanding streaming schemas.
