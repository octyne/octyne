# Chat Completions Schema Status

Snapshot date: 2026-07-11.

The OpenAI-compatible `POST /v1/chat/completions` request implements all 37 current top-level parameters through the compatibility, canonical, and OpenAI provider layers.

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
- seed (deprecated)
- store
- parallel_tool_calls
- safety_identifier
- prompt_cache_key
- max_tokens (deprecated)
- user (deprecated)
- prompt_cache_retention (deprecated)
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
- functions (deprecated)
- function_call (deprecated)

The request schema includes role-specific developer, system, user, assistant, tool, and deprecated function messages; text, image, input-audio, file, and refusal content; assistant audio references; function and custom tool calls; tool results; typed response formats; typed streaming and prompt-cache options; and explicit assistant null content.

## Remaining schema work

- Complete non-streaming response messages, finish reasons, annotations, audio, tool calls, and deprecated function calls.
- Replace response-side `any` fields with typed usage and content/refusal log probability structures.
- Complete streaming deltas for refusals, audio, annotations, function/custom tool-call fragments, multiple choices, usage-only chunks, and obfuscation.
- Preserve successful `[DONE]` behavior and cancellation semantics while expanding streaming schemas.
