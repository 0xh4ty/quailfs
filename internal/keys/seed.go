package keys

import (
	"github.com/tyler-smith/go-bip39"
)

func DeriveSeed(mnemonic string) []byte {
	// Generate a Bip39 seed for the mnemonic
	seed := bip39.NewSeed(mnemonic, "")
	return seed
}
