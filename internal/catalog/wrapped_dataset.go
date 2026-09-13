package catalog

import (
	"crypto/rand"
	"github.com/0xh4ty/quailfs/internal/keys"
	"golang.org/x/crypto/chacha20poly1305"
)

type WrappedDataset struct {
	UserID           []byte
	DatasetID        []byte
	Label            string
	UserX25519Pubkey []byte
	EphX25519Pubkey  []byte
	EncDatasetKey    []byte
}

func CreateWrappedDataset(userID []byte, datasetID []byte, datasetKey []byte, label string, userX25519pubkey []byte) (WrappedDataset, error) {
	var wrappedDataset WrappedDataset
	wrappedDataset.UserID = userID
	wrappedDataset.DatasetID = datasetID
	wrappedDataset.Label = label
	wrappedDataset.UserX25519Pubkey = userX25519pubkey

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return WrappedDataset{}, err
	}

	EphX25519Privkey, EphX25519Pubkey, err := keys.DeriveX25519Keypair(seed)
	if err != nil {
		return WrappedDataset{}, err
	}

	wrappedDataset.EphX25519Pubkey = EphX25519Pubkey

	secretkey, err := keys.ComputeSharedSecret(EphX25519Privkey, userX25519pubkey)
	if err != nil {
		return WrappedDataset{}, err
	}

	aead, err := chacha20poly1305.NewX(secretkey)
	if err != nil {
		return WrappedDataset{}, err
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(datasetKey)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return WrappedDataset{}, err
	}

	encryptedDatasetKey := aead.Seal(nonce, nonce, datasetKey, nil)
	wrappedDataset.EncDatasetKey = encryptedDatasetKey

	return wrappedDataset, nil
}
