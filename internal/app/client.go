package app

import (
	"fmt"
	"github.com/0xh4ty/quailfs/internal/keys"
)

func RunClient() {
	test_mnemonic := "fresh year gaze process tag evidence feature brown involve morning ripple taste secret mesh slot morning script father divorce box maximum language inject olympic"
	mnemonic_seed := keys.DeriveSeed(test_mnemonic)
	fmt.Println("Mnemonic seed derived.")

	ed25519Seed, _ := keys.DeriveEd25519Seed(mnemonic_seed)
	fmt.Println("Ed25519 seed derived.")

	x25519Seed, _ := keys.DeriveX25519Seed(mnemonic_seed)
	fmt.Println("X25519 seed derived.")

	keys.DeriveEd25519Keypair(ed25519Seed)
	fmt.Println("Ed25519 keypair derived.")

	keys.DeriveX25519Keypair(x25519Seed)
	fmt.Println("X25519 keypair derived.")
}
