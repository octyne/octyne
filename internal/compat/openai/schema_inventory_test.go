package openai

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChatCompletionRequestTopLevelFieldInventory(t *testing.T) {
	want := []string{
		"audio", "frequency_penalty", "function_call", "functions", "logit_bias",
		"logprobs", "max_completion_tokens", "max_tokens", "messages", "metadata",
		"modalities", "model", "moderation", "n", "parallel_tool_calls", "prediction",
		"presence_penalty", "prompt_cache_key", "prompt_cache_options",
		"prompt_cache_retention", "reasoning_effort", "response_format",
		"safety_identifier", "seed", "service_tier", "stop", "store", "stream",
		"stream_options", "temperature", "tool_choice", "tools", "top_logprobs",
		"top_p", "user", "verbosity", "web_search_options",
	}

	assertJSONFieldInventory(t, reflect.TypeOf(ChatCompletionRequest{}), want)
}

func assertJSONFieldInventory(t *testing.T, typ reflect.Type, want []string) {
	t.Helper()
	got := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		name := strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]
		if name != "" && name != "-" {
			got = append(got, name)
		}
	}
	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("JSON fields = %v, want %v", got, want)
	}
}
