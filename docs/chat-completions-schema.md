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

## Completed non-streaming response schema

- Top-level `id`, `object`, `created`, `model`, `choices`, `moderation`,
  `service_tier`, deprecated `system_fingerprint`, and `usage` fields.
- Nullable assistant content and refusal output.
- URL-citation annotations and generated audio data.
- Function and custom tool calls plus deprecated `function_call`.
- Typed `stop`, `length`, `tool_calls`, `content_filter`, and deprecated
  `function_call` finish reasons.
- Typed content and refusal token log probabilities, nullable UTF-8 byte arrays,
  and top-token alternatives.
- Prompt, completion, and total token counts; cached, cache-write, audio,
  reasoning, accepted-prediction, and rejected-prediction token details.
- Typed input and output moderation success/error outcomes.

## Completed streaming response schema

Coverage follows the current
[OpenAI Chat Completions streaming-event reference](https://developers.openai.com/api/reference/resources/chat/subresources/completions/streaming-events#chat.completion.chunk).

- Top-level chunk identity, choices, moderation, service tier, deprecated system
  fingerprint, usage, and stream obfuscation.
- Content, refusal, deprecated function-call, and indexed function tool-call
  fragments.
- Multiple simultaneous choices, typed finish reasons, and typed content and
  refusal log probabilities.
- Preservation of omitted usage, explicit `usage: null` intermediate chunks,
  and the final empty-choice usage-only chunk.
- Successful `[DONE]` framing, cancellation, and single-owner response-body and
  channel closure behavior.

The current official chunk schema does not define streaming audio, annotation,
or custom tool-call fragment fields. Recheck the upstream reference before
adding those shapes if OpenAI expands the streaming schema.

## Completed protocol compatibility

Request and success-response schema coverage is complete for this snapshot.
Chat Completions requires an Octyne client key through the bearer
`Authorization` header and returns an OpenAI-compatible authentication error
before request decoding when the credential is missing or invalid. It also
returns OpenAI-compatible error envelopes and status codes, adds an
Octyne-generated `x-request-id` to every response, forwards that value upstream
as `X-Client-Request-Id`, sanitizes upstream server failures, and retains the
provider request ID as internal diagnostic metadata. Streaming setup failures
use normal JSON errors; failures after SSE headers are committed use an error
envelope in a `data:` event and do not send `[DONE]`.

The field inventories and translation tests cover all 37 current top-level
request parameters, typed nested message/tool/content shapes, non-streaming
response fields, streaming chunk fields, and explicit omission/null behavior.
Protocol tests cover authentication, validation, unknown models, upstream rate
limits, safe server errors, request-ID correlation, and mid-stream failures.

No Chat Completions request, success-response, error-envelope, request-ID, or
stream-framing compatibility item is known to be missing as of this snapshot.
Future OpenAI schema changes still require a fresh reference audit.
