package keys

import (
	"github.com/tyler-smith/go-bip39"
)

func GenerateMnemonic() string {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	return mnemonic
}
