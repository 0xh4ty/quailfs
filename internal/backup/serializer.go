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
