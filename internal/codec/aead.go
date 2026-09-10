package codec

import (
	"crypto/aes"
	"crypto/cipher"
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
