package keys

import (
	"testing"
)

func TestDeriveEd25519Seed(t *testing.T) {
	test_mnemonic := "fresh year gaze process tag evidence feature brown involve morning ripple taste secret mesh slot morning script father divorce box maximum language inject olympic"
	mnemonic_seed := DeriveSeed(test_mnemonic)

	seed_1, _ := DeriveEd25519Seed(mnemonic_seed)
	seed_2, _ := DeriveEd25519Seed(mnemonic_seed)

	for index := range seed_1 {
		if seed_1[index] != seed_2[index] {
			t.Error("DeriveEd25519Seed: seed_1 and seed_2 are not equal")
		}
	}
}

func TestDeriveX25519Seed(t *testing.T) {
	test_mnemonic := "fresh year gaze process tag evidence feature brown involve morning ripple taste secret mesh slot morning script father divorce box maximum language inject olympic"
	mnemonic_seed := DeriveSeed(test_mnemonic)

	seed_1, _ := DeriveX25519Seed(mnemonic_seed)
	seed_2, _ := DeriveX25519Seed(mnemonic_seed)

	for index := range seed_1 {
		if seed_1[index] != seed_2[index] {
			t.Error("DeriveX25519Seed: seed_1 and seed_2 are not equal")
		}
	}
}
