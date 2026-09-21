package pipeline

import (
	"fmt"
	"github.com/0xh4ty/quailfs/internal/backup"
	"github.com/0xh4ty/quailfs/internal/catalog"
	"github.com/0xh4ty/quailfs/internal/codec"
	"github.com/0xh4ty/quailfs/internal/keys"
	"github.com/0xh4ty/quailfs/pkg/types"
	"os"
	"path/filepath"
)

type BackupResult struct {
	Manifest       types.ManifestEnvelope
	Head           types.Head
	WrappedDataset types.WrappedDataset
	UserIndex      types.UserIndex

	StripeIDs       [][]byte
	ShardCollection [][][]byte
}

type backupShards struct {
	stripeIDs [][]byte
	shards    [][][]byte
}

func backupFile(
	path string,
	datasetID []byte,
	datasetKey []byte,
) (types.File, []types.MChunk, []types.MStripe, backupShards, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.File{}, nil, nil, backupShards{},
			fmt.Errorf("read %q: %w", path, err)
	}

	qchunks, err := backup.Chunker(data)
	if err != nil {
		return types.File{}, nil, nil, backupShards{},
			fmt.Errorf("chunk %q: %w", path, err)
	}

	compressedQChunks, err := codec.Compressor(qchunks)
	if err != nil {
		return types.File{}, nil, nil, backupShards{},
			fmt.Errorf("compress %q: %w", path, err)
	}

	packedStripePlain := backup.PackChunks(compressedQChunks)
	packedStripeSerialized := backup.SerializePackedPlainCollection(
		packedStripePlain,
	)

	file := types.File{
		FileName: filepath.Base(path),
	}

	var mChunks []types.MChunk
	var mStripes []types.MStripe
	var stripeIDs [][]byte
	var shardCollection [][][]byte

	for _, qchunk := range qchunks {
		mChunks = append(mChunks, types.MChunk{
			ChunkID:  qchunk.ChunkID,
			Size:     uint64(len(qchunk.Data)),
			StripeID: nil,
		})

		file.ChunkIDs = append(file.ChunkIDs, qchunk.ChunkID)
	}

	for _, stripe := range packedStripeSerialized {
		stripeID := keys.DeriveStripeID(datasetID, stripe)

		stripeKey, err := keys.DeriveStripeKey(
			datasetKey,
			datasetID,
			stripeID,
		)
		if err != nil {
			return types.File{}, nil, nil, backupShards{},
				fmt.Errorf("derive stripe key: %w", err)
		}

		nonce96, err := keys.DeriveNonce96(
			stripeKey,
			datasetID,
			stripeID,
		)
		if err != nil {
			return types.File{}, nil, nil, backupShards{},
				fmt.Errorf("derive nonce: %w", err)
		}

		encryptedStripe, err := codec.EncryptStripe(
			stripeKey,
			nonce96,
			stripe,
		)
		if err != nil {
			return types.File{}, nil, nil, backupShards{},
				fmt.Errorf("encrypt stripe: %w", err)
		}

		shards, err := codec.EncodeStripe(encryptedStripe)
		if err != nil {
			return types.File{}, nil, nil, backupShards{},
				fmt.Errorf("encode stripe: %w", err)
		}

		stripeIDs = append(stripeIDs, stripeID)
		shardCollection = append(shardCollection, shards)

		mStripes = append(mStripes, types.MStripe{
			StripeID:   stripeID,
			K:          8,
			N:          12,
			PayloadLen: uint64(len(encryptedStripe)),
			ShardNames: nil,
		})
	}

	return file, mChunks, mStripes, backupShards{
		stripeIDs: stripeIDs,
		shards:    shardCollection,
	}, nil
}

func Backup(
	paths []string,
	userID []byte,
	datasetID []byte,
	datasetKey []byte,
	catalogKey []byte,
	label string,
	generation uint64,
	parentManifestID []byte,
	userX25519Pubkey []byte,
	ed25519PrivateKey []byte,
) error {
	var stripeIDs [][]byte
	var shardCollection [][][]byte

	var trees []types.Tree
	var mChunks []types.MChunk
	var mStripes []types.MStripe

	for _, rootPath := range paths {
		info, err := os.Stat(rootPath)
		if err != nil {
			return fmt.Errorf("stat %q: %w", rootPath, err)
		}

		tree := types.Tree{
			RootDirectory: filepath.Base(rootPath),
		}

		if info.IsDir() {
			err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				file, chunks, stripes, shards, err := backupFile(
					path,
					datasetID,
					datasetKey,
				)
				if err != nil {
					return err
				}

				tree.Files = append(tree.Files, file)
				mChunks = append(mChunks, chunks...)
				mStripes = append(mStripes, stripes...)
				stripeIDs = append(stripeIDs, shards.stripeIDs...)
				shardCollection = append(shardCollection, shards.shards...)

				return nil
			})
			if err != nil {
				return fmt.Errorf("backup directory %q: %w", rootPath, err)
			}
		} else {
			file, chunks, stripes, shards, err := backupFile(
				rootPath,
				datasetID,
				datasetKey,
			)
			if err != nil {
				return err
			}

			tree.Files = append(tree.Files, file)
			mChunks = append(mChunks, chunks...)
			mStripes = append(mStripes, stripes...)
			stripeIDs = append(stripeIDs, shards.stripeIDs...)
			shardCollection = append(shardCollection, shards.shards...)
		}

		trees = append(trees, tree)
	}

	manifest, err := catalog.CreateManifest(
		datasetID,
		generation,
		parentManifestID,
		trees,
		mChunks,
		mStripes,
		nil,
		catalogKey,
		ed25519PrivateKey,
	)
	if err != nil {
		return fmt.Errorf("create manifest: %w", err)
	}

	_, err = catalog.CreateWrappedDataset(
		userID,
		datasetID,
		datasetKey,
		label,
		generation,
		userX25519Pubkey,
		ed25519PrivateKey,
	)
	if err != nil {
		return fmt.Errorf("create wrapped dataset: %w", err)
	}

	_ = catalog.CreateHead(
		1,
		userID,
		datasetID,
		generation,
		manifest.ManifestID,
		nil,
		ed25519PrivateKey,
	)

	datasetEntry := types.DatasetEntry{
		DatasetID:  datasetID,
		Label:      label,
		Generation: generation,
	}

	_ = catalog.CreateUserIndex(
		1,
		userID,
		generation,
		[]types.DatasetEntry{datasetEntry},
		ed25519PrivateKey,
	)

	_ = stripeIDs
	_ = shardCollection

	return nil
}
