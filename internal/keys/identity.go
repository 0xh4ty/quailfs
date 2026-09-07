package keys

import (
	"bytes"
	"crypto/sha256"
)

func DeriveUserID(ed25519publicKey []byte) [32]byte {
	var buf bytes.Buffer
	buf.Write([]byte("quailfs/user/v1"))
	buf.Write(ed25519publicKey)
	hash := sha256.Sum256(buf.Bytes())
	return hash
}
