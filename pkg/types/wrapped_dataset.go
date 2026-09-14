package types

type WrappedDataset struct {
	Body WrappedDatasetBody
	Sig  Sig
}

type WrappedDatasetBody struct {
	UserID           []byte
	DatasetID        []byte
	Label            string
	Generation       uint64
	UserX25519Pubkey []byte
	EphX25519Pubkey  []byte
	EncDatasetKey    []byte
}

type Sig struct {
	Signature []byte
}
