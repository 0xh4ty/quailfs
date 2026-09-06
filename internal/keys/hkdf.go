package keys

import (
	"crypto/hkdf"
	"crypto/sha256"
)

func DeriveEd25519Seed(mnemonic_seed []byte) ([]byte, error) {
	ed25519_seed, err := hkdf.Key(sha256.New, mnemonic_seed, []byte("quailfs"), "quailfs/ed25519/v1", 32)

	return ed25519_seed, err
}

func DeriveX25519Seed(mnemonic_seed []byte) ([]byte, error) {
	x25519_seed, err := hkdf.Key(sha256.New, mnemonic_seed, []byte("quailfs"), "quailfs/x25519/v1", 32)

	return x25519_seed, err
}
