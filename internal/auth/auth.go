package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
)

type Verifier interface {
	Verify(ctx context.Context, key string) (bool, error)
}

type StaticKeyVerifier struct {
	keyHashes [][sha256.Size]byte
}

func NewStaticKeyVerifier(keys []string) *StaticKeyVerifier {
	keyHashes := make([][sha256.Size]byte, 0, len(keys))

	for _, key := range keys {
		if key == "" {
			continue
		}

		keyHashes = append(keyHashes, sha256.Sum256([]byte(key)))
	}

	return &StaticKeyVerifier{
		keyHashes: keyHashes,
	}
}

func (v *StaticKeyVerifier) Verify(
	_ context.Context,
	key string,
) (bool, error) {
	candidateHash := sha256.Sum256([]byte(key))
	matched := 0

	for i := range v.keyHashes {
		matched |= subtle.ConstantTimeCompare(
			candidateHash[:],
			v.keyHashes[i][:],
		)
	}

	return matched == 1, nil
}
