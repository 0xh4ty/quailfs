package keys

import (
	"bytes"
	"crypto/hkdf"
	"crypto/sha256"
)

func DeriveStripeID(datasetID []byte, packed_plain []byte) []byte {
	packedHash := sha256.Sum256(packed_plain)
	var buf bytes.Buffer
	buf.Write([]byte("quailfs/stripe/v1"))
	buf.Write(datasetID)
	buf.Write(packedHash[:])
	stripeID := sha256.Sum256(buf.Bytes())
	return stripeID[:]
}

func DeriveStripeKey(dataKey []byte, datasetID []byte, stripeID []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte("quailfs/data/v1"))
	buf.Write(stripeID)
	stripeKey, err := hkdf.Key(sha256.New, dataKey, datasetID, buf.String(), 32)
	if err != nil {
		return nil, err
	}
	return stripeKey, nil
}

func DeriveNonce96(stripeKey []byte, datasetID []byte, stripeID []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte("quailfs/nonce/v1"))
	buf.Write(stripeID)
	nonce96, err := hkdf.Key(sha256.New, stripeKey, datasetID, buf.String(), 12)
	if err != nil {
		return nil, err
	}
	return nonce96, nil
}
