package codec

import (
	"github.com/0xh4ty/quailfs/internal/backup"
	"github.com/klauspost/compress/zstd"
)

func Compressor(qchunks []backup.QChunk) ([]backup.QChunk, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	for i := 0; i < len(qchunks); i++ {
		qchunk := &qchunks[i]
		qchunk.Data = encoder.EncodeAll(qchunk.Data, nil)
		qchunk.CLength = uint64(len(qchunk.Data))
	}

	return qchunks, nil
}
