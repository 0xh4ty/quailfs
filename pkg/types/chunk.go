package types

type QChunk struct {
	Offset  uint64
	Length  uint64
	CLength uint64
	Data    []byte
	ChunkID []byte
}
