package types

type ManifestEnvelope struct {
	ManifestPlain ManifestPlain
	Envelope      []byte
	Sig           Sig
	ManifestID    []byte
}

type ManifestPlain struct {
	DatasetID        []byte
	Generation       uint64
	ParentManifestID []byte
	Trees            []Tree
	MChunks          []MChunk
	MStripes         []MStripe
	Tombstones       []string
}

type MChunk struct {
	ChunkID  []byte
	Size     uint64
	StripeID []byte
}

type MStripe struct {
	StripeID   []byte
	K          uint64
	N          uint64
	PayloadLen uint64
	ShardNames [][]byte
}

type Tree struct {
	RootDirectory    string
	Files            []File
	ChildDirectories []*Tree
}

type File struct {
	FileName string
	ChunkIDs [][]byte
}
