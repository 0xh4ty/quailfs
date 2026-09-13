package backup

import (
	"bytes"
	"encoding/binary"
	"github.com/0xh4ty/quailfs/pkg/types"
)

func serializeQChunks(qchunks []types.QChunk) []byte {
	count := 0
	var serializedQChunks []byte

	for count < len(qchunks) {
		var buf bytes.Buffer

		offsetBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(offsetBytes, qchunks[count].Offset)
		buf.Write(offsetBytes)

		lengthBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(lengthBytes, qchunks[count].Length)
		buf.Write(lengthBytes)

		clengthBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(clengthBytes, qchunks[count].CLength)
		buf.Write(clengthBytes)

		buf.Write(qchunks[count].Data)
		buf.Write(qchunks[count].ChunkID)

		serializedQChunks = append(serializedQChunks, buf.Bytes()...)

		count++
	}

	return serializedQChunks
}

func SerializePackedPlainCollection(packedPlainCollection []types.PackedPlain) [][]byte {
	count := 0
	var serializedPackedPlainCollection [][]byte

	for count < len(packedPlainCollection) {
		var buf bytes.Buffer
		buf.Write([]byte(packedPlainCollection[count].Version))

		payloadLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(payloadLenBytes, packedPlainCollection[count].PayloadLen)
		buf.Write(payloadLenBytes)

		serializedQChunks := serializeQChunks(packedPlainCollection[count].QChunks)
		buf.Write(serializedQChunks)

		buf.Write(packedPlainCollection[count].Padding)

		serializedPackedPlainCollection = append(serializedPackedPlainCollection, buf.Bytes())

		count++
	}

	return serializedPackedPlainCollection
}

func deserializeQChunks(serializedQChunks []byte) []types.QChunk {
	var qchunks []types.QChunk
	position := 0
	for position < len(serializedQChunks) {
		var qchunk types.QChunk
		qchunk.Offset = binary.BigEndian.Uint64(serializedQChunks[position : position+8])
		position += 8
		qchunk.Length = binary.BigEndian.Uint64(serializedQChunks[position : position+8])
		position += 8
		qchunk.CLength = binary.BigEndian.Uint64(serializedQChunks[position : position+8])
		position += 8
		qchunk.Data = serializedQChunks[position : position+int(qchunk.CLength)]
		position += int(qchunk.CLength)
		qchunk.ChunkID = serializedQChunks[position : position+32]
		position += 32
		qchunks = append(qchunks, qchunk)
	}
	return qchunks
}

func DeserializePackedPlainCollection(serializedPackedPlainCollection [][]byte) []types.PackedPlain {
	var packedPlainCollection []types.PackedPlain
	capacity := 2 * 1024 * 1024

	for i := range len(serializedPackedPlainCollection) {
		position := 0
		serializedPackedPlain := serializedPackedPlainCollection[i]
		for position < len(serializedPackedPlain) {
			var packedPlain types.PackedPlain
			packedPlain.Version = string(serializedPackedPlain[position : position+2])
			position += 2
			packedPlain.PayloadLen = binary.BigEndian.Uint64(serializedPackedPlain[position : position+8])
			position += 8
			qchunksLen := int(packedPlain.PayloadLen) - 2 - 8
			packedPlain.QChunks = deserializeQChunks(serializedPackedPlain[position : position+qchunksLen])
			position += qchunksLen
			paddingSize := capacity - int(packedPlain.PayloadLen)
			packedPlain.Padding = serializedPackedPlain[position : position+paddingSize]
			position += paddingSize

			packedPlainCollection = append(packedPlainCollection, packedPlain)
		}
	}

	return packedPlainCollection
}

func SerializeWrappedDatasetBody(wrappedDataset types.WrappedDataset) []byte {
	var serializedWrappedDatasetBody []byte

	var buf bytes.Buffer
	buf.Write(wrappedDataset.Body.UserID)
	buf.Write(wrappedDataset.Body.DatasetID)
	labelLen := uint64(len([]byte(wrappedDataset.Body.Label)))
	labelLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(labelLenBytes, labelLen)
	buf.Write(labelLenBytes)
	buf.Write([]byte(wrappedDataset.Body.Label))
	buf.Write(wrappedDataset.Body.UserX25519Pubkey)
	buf.Write(wrappedDataset.Body.EphX25519Pubkey)
	EncDatasetKeyLen := uint64(len([]byte(wrappedDataset.Body.EncDatasetKey)))
	EncDatasetKeyLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(EncDatasetKeyLenBytes, EncDatasetKeyLen)
	buf.Write(EncDatasetKeyLenBytes)
	buf.Write(wrappedDataset.Body.EncDatasetKey)

	serializedWrappedDatasetBody = buf.Bytes()

	return serializedWrappedDatasetBody
}

func SerializeDatasetEntries(datasets []types.DatasetEntry) []byte {
	var serializedDatasetEntries []byte
	var buf bytes.Buffer

	datasetEntryCount := uint64(len(datasets))
	datasetEntryCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(datasetEntryCountBytes, datasetEntryCount)
	buf.Write(datasetEntryCountBytes)

	for i := range datasets {
		buf.Write(datasets[i].DatasetID)

		labelLen := uint64(len([]byte(datasets[i].Label)))
		labelLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(labelLenBytes, labelLen)
		buf.Write(labelLenBytes)

		buf.Write([]byte(datasets[i].Label))

		generationBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(generationBytes, datasets[i].Generation)
		buf.Write(generationBytes)
	}

	serializedDatasetEntries = buf.Bytes()

	return serializedDatasetEntries
}

func SerializeUserIndexBody(userIndex types.UserIndex) []byte {
	var serializedUserIndexBody []byte

	var buf bytes.Buffer
	schemaBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(schemaBytes, userIndex.Body.Schema)
	buf.Write(schemaBytes)

	buf.Write(userIndex.Body.UserID)

	generationBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(generationBytes, userIndex.Body.Generation)
	buf.Write(generationBytes)

	serializedDatasetEntries := SerializeDatasetEntries(userIndex.Body.Datasets)
	buf.Write(serializedDatasetEntries)

	buf.Write([]byte(userIndex.Body.UpdatedAt))

	serializedUserIndexBody = buf.Bytes()

	return serializedUserIndexBody
}
