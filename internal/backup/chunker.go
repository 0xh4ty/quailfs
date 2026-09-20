package backup

import (
	"bytes"
	"crypto/sha256"
	"github.com/0xh4ty/quailfs/pkg/types"
	"github.com/jotfs/fastcdc-go"
	"io"
	"sort"
)

func Chunker(data []byte) ([]types.QChunk, error) {
	dataBuf := bytes.NewReader(data)

	opts := fastcdc.Options{
		MinSize:     64 * 1024,
		AverageSize: 256 * 1024,
		MaxSize:     1 * 1024 * 1024,
	}

	chunker, err := fastcdc.NewChunker(dataBuf, opts)
	if err != nil {
		return nil, err
	}

	var qchunks []types.QChunk

	for {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunk_data := append([]byte(nil), chunk.Data...)

		hash := sha256.Sum256(chunk.Data)
		chunkID := hash[:]

		var qchunk types.QChunk
		qchunk.Offset = uint64(chunk.Offset)
		qchunk.Length = uint64(chunk.Length)
		qchunk.Data = chunk_data
		qchunk.ChunkID = chunkID

		qchunks = append(qchunks, qchunk)
	}

	return qchunks, nil
}

func Dechunker(qchunks []types.QChunk) ([]byte, error) {
	var dataBuf bytes.Buffer
	sort.Slice(qchunks, func(i, j int) bool {
		return qchunks[i].Offset < qchunks[j].Offset
	})
	count := 0

	for count < len(qchunks) {
		dataBuf.Write(qchunks[count].Data)
		count++
	}

	data := dataBuf.Bytes()

	return data, nil
}
