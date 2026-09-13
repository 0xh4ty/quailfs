package types

type UserIndex struct {
	Body UserIndexBody
	Sig  Sig
}

type UserIndexBody struct {
	Schema     uint64
	UserID     []byte
	Generation uint64
	Datasets   []DatasetEntry
	UpdatedAt  string
}

type DatasetEntry struct {
	DatasetID  []byte
	Label      string
	Generation uint64
}
