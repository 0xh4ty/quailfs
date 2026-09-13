package types

type WrappedDataset struct {
	Body Body
	Sig  Sig
}

type Body struct {
	UserID           []byte
	DatasetID        []byte
	Label            string
	UserX25519Pubkey []byte
	EphX25519Pubkey  []byte
	EncDatasetKey    []byte
}

type Sig struct {
	Signature []byte
}
