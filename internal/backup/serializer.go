package backup

import (
	"bytes"
	"encoding/binary"
)

func serializeQChunks(qchunks []QChunk) []byte {
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

func SerializePackedPlainCollection(packedPlainCollection []PackedPlain) [][]byte {
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

func deserializeQChunks(serializedQChunks []byte) []QChunk {
	var qchunks []QChunk
	position := 0
	for position < len(serializedQChunks) {
		var qchunk QChunk
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

func DeserializePackedPlainCollection(serializedPackedPlainCollection [][]byte) []PackedPlain {
	var packedPlainCollection []PackedPlain
	capacity := 2 * 1024 * 1024

	for i := range len(serializedPackedPlainCollection) {
		position := 0
		serializedPackedPlain := serializedPackedPlainCollection[i]
		for position < len(serializedPackedPlain) {
			var packedPlain PackedPlain
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
