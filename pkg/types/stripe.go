package types

type PackedPlain struct {
	Version    string
	PayloadLen uint64
	QChunks    []QChunk
	Padding    []byte
}
