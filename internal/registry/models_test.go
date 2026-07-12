package registry

import (
	"reflect"
	"testing"
)

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

func TestRegistryListReturnsSortedSnapshot(t *testing.T) {
	registry := NewRegistry()
	registry.Register("openai/gpt-5-nano", Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	})
	registry.Register("openai/gpt-4.1-mini", Model{
		Provider: "openai",
		ModelID:  "gpt-4.1-mini",
	})

	want := []RegisteredModel{
		{
			Name: "openai/gpt-4.1-mini",
			Model: Model{
				Provider: "openai",
				ModelID:  "gpt-4.1-mini",
			},
		},
		{
			Name: "openai/gpt-5-nano",
			Model: Model{
				Provider: "openai",
				ModelID:  "gpt-5-nano",
			},
		},
	}

	got := registry.List()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("List() = %+v, want %+v", got, want)
	}

	got[0].Name = "changed"
	if after := registry.List(); !reflect.DeepEqual(after, want) {
		t.Errorf("List() after caller mutation = %+v, want %+v", after, want)
	}
}

func TestRegistryListEmpty(t *testing.T) {
	models := NewRegistry().List()
	if models == nil {
		t.Fatal("List() = nil, want empty slice")
	}
	if len(models) != 0 {
		t.Errorf("len(List()) = %d, want 0", len(models))
	}
}
