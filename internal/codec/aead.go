package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func EncryptStripe(stripeKey []byte, nonce_96 []byte, stripe []byte) ([]byte, error) {

	block, err := aes.NewCipher(stripeKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encryptedStripe := gcm.Seal(nil, nonce_96, stripe, nil)

	return encryptedStripe, nil
}

func DecryptStripe(stripeKey []byte, nonce_96 []byte, encryptedStripe []byte) ([]byte, error) {

	block, err := aes.NewCipher(stripeKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	stripe, err := gcm.Open(nil, nonce_96, encryptedStripe, nil)
	if err != nil {
		return nil, err
	}

	return stripe, nil
}

func EncryptManifest(catalogKey []byte, serializedManifestPlain []byte) ([]byte, error) {

	block, err := aes.NewCipher(catalogKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	encryptedManifest := gcm.Seal(nil, nonce, serializedManifestPlain, nil)

	var buf bytes.Buffer
	buf.Write(nonce)
	buf.Write(encryptedManifest)
	encryptedManifestWithNonce := buf.Bytes()
	return encryptedManifestWithNonce, nil
}
