package catalog

import (
	"crypto/ed25519"
	"fmt"
	"github.com/0xh4ty/quailfs/internal/backup"
)

func VerifyCatalogObject(objectType uint8, data []byte, publicKey []byte) error {
	if len(publicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid Ed25519 public key length: %d", len(publicKey))
	}

	switch objectType {
	case 1:
		userIndex, err := backup.DeserializeUserIndex(data)
		if err != nil {
			return fmt.Errorf("deserialize user index: %w", err)
		}

		body := backup.SerializeUserIndexBody(userIndex)

		message := append([]byte("quailfs/userindex/v1"), body...)

		if !ed25519.Verify(ed25519.PublicKey(publicKey), message, userIndex.Sig.Signature) {
			return fmt.Errorf("invalid user index signature")
		}

		return nil

	case 2:
		head, wrappedDataset, err := backup.DeserializeHeadCatalog(data)
		if err != nil {
			return fmt.Errorf("deserialize head catalog: %w", err)
		}

		headBody := backup.SerializeHeadBody(head)
		headMessage := append([]byte("quailfs/head/v1"), headBody...)

		if !ed25519.Verify(ed25519.PublicKey(publicKey), headMessage, head.Sig.Signature) {
			return fmt.Errorf("invalid head signature")
		}

		wrappedDatasetBody := backup.SerializeWrappedDatasetBody(wrappedDataset)
		wrappedDatasetMessage := append([]byte("quailfs/wrapped/v1"), wrappedDatasetBody...)

		if !ed25519.Verify(ed25519.PublicKey(publicKey), wrappedDatasetMessage, wrappedDataset.Sig.Signature) {
			return fmt.Errorf("invalid wrapped dataset signature")
		}

		return nil

	case 4:
		manifest, err := backup.DeserializeManifestEnvelope(data)
		if err != nil {
			return fmt.Errorf("deserialize manifest envelope: %w", err)
		}

		message := append([]byte("quailfs/manifest/v1"), manifest.Envelope...)

		if !ed25519.Verify(ed25519.PublicKey(publicKey), message, manifest.Sig.Signature) {
			return fmt.Errorf("invalid manifest signature")
		}

		return nil

	default:
		return fmt.Errorf("unknown catalog object type: %d", objectType)
	}
}
