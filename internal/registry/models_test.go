package registry

import "testing"

func TestRegistryRegisterAndGet(t *testing.T) {
	registry := NewRegistry()
	want := Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	}

	registry.Register("openai/gpt-5-nano", want)

	got, ok := registry.Get("openai/gpt-5-nano")
	if !ok {
		t.Fatal("Get(\"openai/gpt-5-nano\") ok = false, want true")
	}
	if got != want {
		t.Errorf("Get(\"openai/gpt-5-nano\") = %+v, want %+v", got, want)
	}
}

func TestRegistryInstancesAreIndependent(t *testing.T) {
	first := NewRegistry()
	second := NewRegistry()

	first.Register("openai/gpt-5-nano", Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	})

	if _, ok := second.Get("openai/gpt-5-nano"); ok {
		t.Error("second registry contains model registered only in first registry")
	}
}
