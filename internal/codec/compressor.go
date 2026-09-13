package codec

import (
	"github.com/0xh4ty/quailfs/pkg/types"
	"github.com/klauspost/compress/zstd"
)

func Compressor(qchunks []types.QChunk) ([]types.QChunk, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	for i := range qchunks {
		qchunk := &qchunks[i]
		qchunk.Data = encoder.EncodeAll(qchunk.Data, nil)
		qchunk.CLength = uint64(len(qchunk.Data))
	}

	return qchunks, nil
}

func Decompressor(qchunks []types.QChunk) ([]types.QChunk, error) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	for i := range qchunks {
		qchunk := &qchunks[i]
		qchunk.Data, err = decoder.DecodeAll(qchunk.Data, nil)
		if err != nil {
			return nil, err
		}
	}

	return qchunks, nil
}
