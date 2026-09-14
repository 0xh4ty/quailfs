package catalog

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"github.com/0xh4ty/quailfs/internal/backup"
	"github.com/0xh4ty/quailfs/internal/keys"
	"github.com/0xh4ty/quailfs/pkg/types"
	"golang.org/x/crypto/chacha20poly1305"
)

func CreateWrappedDataset(userID []byte, datasetID []byte, datasetKey []byte, label string, generation uint64, userX25519pubkey []byte, ed25519privateKey []byte) (types.WrappedDataset, error) {
	var wrappedDataset types.WrappedDataset
	wrappedDataset.Body.UserID = userID
	wrappedDataset.Body.DatasetID = datasetID
	wrappedDataset.Body.Label = label
	wrappedDataset.Body.Generation = generation
	wrappedDataset.Body.UserX25519Pubkey = userX25519pubkey

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return types.WrappedDataset{}, err
	}

	EphX25519Privkey, EphX25519Pubkey, err := keys.DeriveX25519Keypair(seed)
	if err != nil {
		return types.WrappedDataset{}, err
	}

	wrappedDataset.Body.EphX25519Pubkey = EphX25519Pubkey

	secretkey, err := keys.ComputeSharedSecret(EphX25519Privkey, userX25519pubkey)
	if err != nil {
		return types.WrappedDataset{}, err
	}

	aead, err := chacha20poly1305.NewX(secretkey)
	if err != nil {
		return types.WrappedDataset{}, err
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(datasetKey)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return types.WrappedDataset{}, err
	}

	encryptedDatasetKey := aead.Seal(nonce, nonce, datasetKey, nil)
	wrappedDataset.Body.EncDatasetKey = encryptedDatasetKey

	serializedWrappedDatasetBody := backup.SerializeWrappedDatasetBody(wrappedDataset)
	var buf bytes.Buffer
	buf.WriteString("quailfs/wrapped/v1")
	buf.Write(serializedWrappedDatasetBody)

	signature := ed25519.Sign(ed25519privateKey, buf.Bytes())

	wrappedDataset.Sig.Signature = signature

	return wrappedDataset, nil
}
