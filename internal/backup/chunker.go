package backup

import (
	"crypto/sha256"
	"github.com/jotfs/fastcdc-go"
	"io"
	"os"
)

type QChunk struct {
	Offset  uint64
	Length  uint64
	CLength uint64
	Data    []byte
	ChunkID []byte
}

func Chunker(file *os.File) ([]QChunk, error) {
	opts := fastcdc.Options{
		MinSize:     64 * 1024,
		AverageSize: 256 * 1024,
		MaxSize:     1 * 1024 * 1024,
	}

	chunker, err := fastcdc.NewChunker(file, opts)
	if err != nil {
		return nil, err
	}

	var qchunks []QChunk

	for {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256(chunk.Data)
		chunkID := hash[:]

		var qchunk QChunk
		qchunk.Offset, qchunk.Length, qchunk.Data, qchunk.ChunkID = uint64(chunk.Offset), uint64(chunk.Length), chunk.Data, chunkID

		qchunks = append(qchunks, qchunk)
	}

	return qchunks, nil
}
