package keys

import (
	"crypto/ecdh"
	"crypto/ed25519"
)

func DeriveEd25519Keypair(ed25519seed []byte) ([]byte, []byte) {
	ed25519privateKey := ed25519.NewKeyFromSeed(ed25519seed)
	ed25519publicKey := ed25519privateKey.Public().(ed25519.PublicKey)
	return ed25519privateKey, ed25519publicKey
}

func DeriveX25519Keypair(x25519seed []byte) ([]byte, []byte, error) {
	x25519privateKey, err := ecdh.X25519().NewPrivateKey(x25519seed)
	if err != nil {
		return nil, nil, err
	}
	x25519publicKey := x25519privateKey.PublicKey()
	return x25519privateKey.Bytes(), x25519publicKey.Bytes(), nil
}
