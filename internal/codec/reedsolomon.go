package codec

import (
	"bytes"
	"github.com/klauspost/reedsolomon"
)

func EncodeStripe(encryptedStripe []byte) ([][]byte, error) {
	enc, err := reedsolomon.New(8, 4)
	if err != nil {
		return nil, err
	}

	split, err := enc.Split(encryptedStripe)
	if err != nil {
		return nil, err
	}
	enc.Encode(split)

	ok, err := enc.Verify(split)
	if ok != true {
		return nil, err
	}

	return split, nil
}

func DecodeStripe(shards [][]byte) ([]byte, error) {
	stripeLen := (2 * 1024 * 1024) + 16
	var buf bytes.Buffer
	enc, err := reedsolomon.New(8, 4)
	if err != nil {
		return nil, err
	}

	err = enc.ReconstructData(shards)
	if err != nil {
		return nil, err
	}

	err = enc.Join(&buf, shards, stripeLen)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
