package auth

import (
	"context"
	"testing"
)

func TestStaticKeyVerifierVerify(t *testing.T) {
	verifier := NewStaticKeyVerifier([]string{
		"first-client-key",
		"second-client-key",
	})

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "first configured key",
			key:  "first-client-key",
			want: true,
		},
		{
			name: "second configured key",
			key:  "second-client-key",
			want: true,
		},
		{
			name: "unknown key",
			key:  "unknown-client-key",
			want: false,
		},
		{
			name: "configured key prefix",
			key:  "first-client",
			want: false,
		},
		{
			name: "configured key suffix",
			key:  "client-key",
			want: false,
		},
		{
			name: "empty key",
			key:  "",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := verifier.Verify(context.Background(), test.key)
			if err != nil {
				t.Fatalf("Verify() error = %v", err)
			}
			if got != test.want {
				t.Errorf("Verify() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestStaticKeyVerifierFailsClosedWithoutConfiguredKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []string
	}{
		{
			name: "nil keys",
			keys: nil,
		},
		{
			name: "empty keys",
			keys: []string{},
		},
		{
			name: "empty configured values",
			keys: []string{"", ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := NewStaticKeyVerifier(test.keys)

			for _, key := range []string{"", "any-client-key"} {
				got, err := verifier.Verify(context.Background(), key)
				if err != nil {
					t.Fatalf("Verify(%q) error = %v", key, err)
				}
				if got {
					t.Errorf("Verify(%q) = true, want false", key)
				}
			}
		})
	}
}

func TestStaticKeyVerifierIgnoresEmptyConfiguredValues(t *testing.T) {
	verifier := NewStaticKeyVerifier([]string{"", "configured-client-key", ""})

	got, err := verifier.Verify(context.Background(), "configured-client-key")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if !got {
		t.Error("Verify() = false, want true")
	}

	got, err = verifier.Verify(context.Background(), "")
	if err != nil {
		t.Fatalf("Verify(empty key) error = %v", err)
	}
	if got {
		t.Error("Verify(empty key) = true, want false")
	}
}

func TestStaticKeyVerifierImplementsVerifier(t *testing.T) {
	var verifier Verifier = NewStaticKeyVerifier([]string{"configured-client-key"})

	if verifier == nil {
		t.Fatal("Verifier = nil, want static key verifier")
	}
}
