package catalog

import (
	"bytes"
	"crypto/ed25519"
	"github.com/0xh4ty/quailfs/internal/backup"
	"github.com/0xh4ty/quailfs/pkg/types"
	"time"
)

func CreateUserIndex(schema uint64, userID []byte, generation uint64, datasetEntries []types.DatasetEntry, ed25519privateKey []byte) types.UserIndex {
	var userIndex types.UserIndex

	userIndex.Body.Schema = schema
	userIndex.Body.UserID = userID
	userIndex.Body.Generation = generation
	userIndex.Body.Datasets = datasetEntries

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	userIndex.Body.UpdatedAt = updatedAt

	serializeUserIndexBody := backup.SerializeUserIndexBody(userIndex)
	var buf bytes.Buffer
	buf.WriteString("quailfs/userindex/v1")
	buf.Write(serializeUserIndexBody)

	signature := ed25519.Sign(ed25519privateKey, buf.Bytes())
	userIndex.Sig.Signature = signature

	return userIndex
}
