package openai

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChatCompletionResponseTopLevelFieldInventory(t *testing.T) {
	want := []string{
		"choices", "created", "id", "model", "moderation", "object", "service_tier",
		"system_fingerprint", "usage",
	}
	assertTopLevelJSONFields(t, reflect.TypeOf(ChatCompletionResponse{}), want)
}

func TestChatCompletionChunkTopLevelFieldInventory(t *testing.T) {
	want := []string{
		"choices", "created", "id", "model", "moderation", "object", "obfuscation",
		"service_tier", "system_fingerprint",
	}
	assertTopLevelJSONFields(t, reflect.TypeOf(ChatCompletionChunk{}), want)
}

func assertTopLevelJSONFields(t *testing.T, typ reflect.Type, want []string) {
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
