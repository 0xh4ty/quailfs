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

func SerializeCatalogPeerHints(catalogPeerHints [][]byte) []byte {
	var serializedCatalogPeerHints []byte

	var buf bytes.Buffer

	for i := range catalogPeerHints {
		buf.Write(catalogPeerHints[i])
	}

	serializedCatalogPeerHints = buf.Bytes()

	return serializedCatalogPeerHints
}

func SerializeHeadBody(head types.Head) []byte {
	var serializedHeadBody []byte

	var buf bytes.Buffer
	schemaBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(schemaBytes, head.Body.Schema)
	buf.Write(schemaBytes)

	buf.Write(head.Body.UserID)
	buf.Write(head.Body.DatasetID)

	generationBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(generationBytes, head.Body.Generation)
	buf.Write(generationBytes)

	buf.Write(head.Body.ManifestID)

	catalogPeerHintsCount := uint64(len(head.Body.CatalogPeerHints))
	catalogPeerHintsCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(catalogPeerHintsCountBytes, catalogPeerHintsCount)
	buf.Write(catalogPeerHintsCountBytes)

	serializedCatalogPeerHints := SerializeCatalogPeerHints(head.Body.CatalogPeerHints)
	buf.Write(serializedCatalogPeerHints)

	buf.Write([]byte(head.Body.CreatedAt))
	serializedHeadBody = buf.Bytes()

	return serializedHeadBody
}

func SerializeFiles(files []types.File) []byte {
	var serializedFiles []byte

	var buf bytes.Buffer

	filesCount := uint64(len(files))
	filesCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(filesCountBytes, filesCount)
	buf.Write(filesCountBytes)

	for i := range files {
		fileNameLen := uint64(len([]byte(files[i].FileName)))
		fileNameLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(fileNameLenBytes, fileNameLen)
		buf.Write(fileNameLenBytes)

		buf.Write([]byte(files[i].FileName))

		chunkIDCount := uint64(len(files[i].ChunkIDs))
		chunkIDCountBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(chunkIDCountBytes, chunkIDCount)
		buf.Write(chunkIDCountBytes)

		for j := range files[i].ChunkIDs {
			buf.Write(files[i].ChunkIDs[j])
		}
	}

	serializedFiles = buf.Bytes()

	return serializedFiles
}

func SerializeTree(tree types.Tree) []byte {
	var serializedTree []byte

	var buf bytes.Buffer

	RootDirectoryLen := uint64(len([]byte(tree.RootDirectory)))
	RootDirectoryLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(RootDirectoryLenBytes, RootDirectoryLen)
	buf.Write(RootDirectoryLenBytes)

	buf.Write([]byte(tree.RootDirectory))
	serializedFiles := SerializeFiles(tree.Files)
	buf.Write(serializedFiles)

	childDirectoryCount := uint64(len(tree.ChildDirectories))
	childDirectoryCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(childDirectoryCountBytes, childDirectoryCount)
	buf.Write(childDirectoryCountBytes)

	for i := range tree.ChildDirectories {
		serializedChildDirectory := SerializeTree(*tree.ChildDirectories[i])
		buf.Write(serializedChildDirectory)
	}

	serializedTree = buf.Bytes()

	return serializedTree
}

func SerializeMChunks(mchunks []types.MChunk) []byte {
	var serializedMChunks []byte

	var buf bytes.Buffer

	mchunksCount := uint64(len(mchunks))
	mchunksCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(mchunksCountBytes, mchunksCount)
	buf.Write(mchunksCountBytes)

	for i := range mchunks {
		buf.Write(mchunks[i].ChunkID)

		sizeBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(sizeBytes, mchunks[i].Size)
		buf.Write(sizeBytes)

		buf.Write(mchunks[i].StripeID)
	}

	serializedMChunks = buf.Bytes()

	return serializedMChunks
}

func SerializeMStripes(mstripes []types.MStripe) []byte {
	var serializedMStripes []byte

	var buf bytes.Buffer

	mstripesCount := uint64(len(mstripes))
	mstripesCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(mstripesCountBytes, mstripesCount)
	buf.Write(mstripesCountBytes)

	for i := range mstripes {
		buf.Write(mstripes[i].StripeID)

		kBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(kBytes, mstripes[i].K)
		buf.Write(kBytes)

		nBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(nBytes, mstripes[i].N)
		buf.Write(nBytes)

		PayloadLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(PayloadLenBytes, mstripes[i].PayloadLen)
		buf.Write(PayloadLenBytes)

		shardNamesCount := uint64(len(mstripes[i].ShardNames))
		shardNamesCountBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(shardNamesCountBytes, shardNamesCount)
		buf.Write(shardNamesCountBytes)

		for j := range mstripes[i].ShardNames {
			buf.Write(mstripes[i].ShardNames[j])
		}
	}

	serializedMStripes = buf.Bytes()

	return serializedMStripes
}

func SerializeTombstones(tombstones []string) []byte {
	var serializedTombstones []byte

	var buf bytes.Buffer

	tombstonesCount := uint64(len(tombstones))
	tombstonesCountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tombstonesCountBytes, tombstonesCount)
	buf.Write(tombstonesCountBytes)

	for i := range tombstones {
		tombstoneLen := uint64(len([]byte(tombstones[i])))
		tombstoneLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(tombstoneLenBytes, tombstoneLen)
		buf.Write(tombstoneLenBytes)

		buf.Write([]byte(tombstones[i]))
	}

	serializedTombstones = buf.Bytes()

	return serializedTombstones
}

func SerializeManifestPlain(manifest types.ManifestEnvelope) []byte {
	var serializedManifestPlain []byte

	var buf bytes.Buffer
	buf.Write(manifest.ManifestPlain.DatasetID)

	generationBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(generationBytes, manifest.ManifestPlain.Generation)
	buf.Write(generationBytes)

	buf.Write(manifest.ManifestPlain.ParentManifestID)

	serializedTree := SerializeTree(manifest.ManifestPlain.Tree)
	buf.Write(serializedTree)

	serializedMChunks := SerializeMChunks(manifest.ManifestPlain.MChunks)
	buf.Write(serializedMChunks)

	serializedMStripes := SerializeMStripes(manifest.ManifestPlain.MStripes)
	buf.Write(serializedMStripes)

	serializedTombstones := SerializeTombstones(manifest.ManifestPlain.Tombstones)
	buf.Write(serializedTombstones)

	serializedManifestPlain = buf.Bytes()

	return serializedManifestPlain
}
