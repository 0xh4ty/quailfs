package backup

import (
	"github.com/0xh4ty/quailfs/pkg/types"
)

func calculatePackedPlainSize(packedPlain types.PackedPlain) uint64 {
	payloadlensize := 8
	var allQchunksSize uint64
	for _, qchunk := range packedPlain.QChunks {
		allQchunksSize += calculateQChunkSize(qchunk)
	}
	packedPlainSize := uint64(len([]byte(packedPlain.Version))) + uint64(payloadlensize) + allQchunksSize
	return packedPlainSize
}

func calculateQChunkSize(qchunk types.QChunk) uint64 {
	var qchunkSize uint64
	offsetsize := 8
	lengthsize := 8
	clengthsize := 8
	qchunkSize = uint64(offsetsize + lengthsize + clengthsize + len(qchunk.Data) + len(qchunk.ChunkID))
	return qchunkSize
}

func PackChunks(qchunks []types.QChunk) []types.PackedPlain {
	var capacity uint64
	capacity = 2 * 1024 * 1024
	var packedPlainCollection []types.PackedPlain
	count := 0
	for {
		if count == len(qchunks) {
			break
		}

		var packedPlain types.PackedPlain
		packedPlain.Version = "v1"

		for count < len(qchunks) {
			if calculatePackedPlainSize(packedPlain)+calculateQChunkSize(qchunks[count]) > capacity {
				break
			}

			packedPlain.QChunks = append(packedPlain.QChunks, qchunks[count])
			count++
		}

		packedPlain.PayloadLen = calculatePackedPlainSize(packedPlain)
		paddingSize := capacity - packedPlain.PayloadLen
		padding := make([]byte, paddingSize)
		packedPlain.Padding = padding
		packedPlainCollection = append(packedPlainCollection, packedPlain)
	}

	return packedPlainCollection
}

func UnpackChunks(packedPlainCollection []types.PackedPlain) []types.QChunk {
	var compressedQChunks []types.QChunk
	for i := range packedPlainCollection {
		packedPlain := packedPlainCollection[i]
		for j := range packedPlain.QChunks {
			compressedQChunks = append(compressedQChunks, packedPlain.QChunks[j])
		}
	}
	return compressedQChunks
}
