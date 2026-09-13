package keys

import (
    "golang.org/x/crypto/curve25519"
)

func ComputeSharedSecret (privateKey []byte, publicKey []byte) ([]byte, error) {
    sharedSecret, err := curve25519.X25519(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

    return sharedSecret, nil
}
