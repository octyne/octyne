package registry

import "testing"

func TestRegistryRegisterAndGet(t *testing.T) {
	registry := NewRegistry()
	want := Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	}

	registry.Register("fast", want)

	got, ok := registry.Get("fast")
	if !ok {
		t.Fatal("Get(\"fast\") ok = false, want true")
	}
	if got != want {
		t.Errorf("Get(\"fast\") = %+v, want %+v", got, want)
	}
}

func TestRegistryInstancesAreIndependent(t *testing.T) {
	first := NewRegistry()
	second := NewRegistry()

	first.Register("fast", Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	})

	if _, ok := second.Get("fast"); ok {
		t.Error("second registry contains model registered only in first registry")
	}
}
