package codec

import (
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
