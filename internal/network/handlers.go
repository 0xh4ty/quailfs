package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	cat "github.com/0xh4ty/quailfs/internal/catalog"
	"github.com/0xh4ty/quailfs/internal/storage"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"go.etcd.io/bbolt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func RegisterHandlers(h host.Host, db *bbolt.DB) {
	h.SetStreamHandler(PutShardProtocol, func(stream network.Stream) {
		handlePutShard(stream, db)
	})

	h.SetStreamHandler(PutCatalogProtocol, func(stream network.Stream) {
		handlePutCatalog(stream, db)
	})
}

func handlePutShard(stream network.Stream, db *bbolt.DB) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)

	shardName, err := readString(reader)
	if err != nil {
		writeResponse(writer, false)
		return
	}

	var shardSize uint64

	if err := binary.Read(reader, binary.BigEndian, &shardSize); err != nil {
		writeResponse(writer, false)
		return
	}

	if shardSize > uint64(^uint(0)>>1) {
		writeResponse(writer, false)
		return
	}

	shard := make([]byte, int(shardSize))

	if _, err := io.ReadFull(reader, shard); err != nil {
		writeResponse(writer, false)
		return
	}

	if !validShardName(shardName) {
		writeResponse(writer, false)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		writeResponse(writer, false)
		return
	}

	shardsDir := filepath.Join(homeDir, ".node-data", "shards")

	if err := os.MkdirAll(shardsDir, 0700); err != nil {
		writeResponse(writer, false)
		return
	}

	shardPath := filepath.Join(shardsDir, shardName)

	if err := os.WriteFile(shardPath, shard, 0600); err != nil {
		writeResponse(writer, false)
		return
	}

	if err := storage.Put(db, "inventory", []byte(shardName), []byte(shardPath)); err != nil {
		os.Remove(shardPath)
		writeResponse(writer, false)
		return
	}

	writeResponse(writer, true)
}

func handlePutCatalog(stream network.Stream, db *bbolt.DB) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)

	objectName, err := readString(reader)
	if err != nil {
		writeResponse(writer, false)
		return
	}

	objectType, err := reader.ReadByte()
	if err != nil {
		writeResponse(writer, false)
		return
	}

	ed25519PublicKey := make([]byte, 32)

	if _, err := io.ReadFull(reader, ed25519PublicKey); err != nil {
		writeResponse(writer, false)
		return
	}

	var catalogSize uint64

	if err := binary.Read(reader, binary.BigEndian, &catalogSize); err != nil {
		writeResponse(writer, false)
		return
	}

	if catalogSize > uint64(^uint(0)>>1) {
		writeResponse(writer, false)
		return
	}

	catalog := make([]byte, int(catalogSize))

	if _, err := io.ReadFull(reader, catalog); err != nil {
		writeResponse(writer, false)
		return
	}

	if err := validateCatalogObjectName(objectName); err != nil {
		writeResponse(writer, false)
		return
	}

	if err := cat.VerifyCatalogObject(objectType, catalog, ed25519PublicKey); err != nil {
		writeResponse(writer, false)
		return
	}

	if err := storage.Put(db, "catalogs", []byte(objectName), catalog); err != nil {
		writeResponse(writer, false)
		return
	}

	writeResponse(writer, true)
}

func validShardName(name string) bool {
	if len(name) != 66 {
		return false
	}

	for _, c := range name {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}

func validateCatalogObjectName(name string) error {
	if name == "" {
		return fmt.Errorf("empty catalog object name")
	}

	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("invalid catalog object name")
	}

	for _, c := range name {
		if c < 0x20 || c == 0x7f {
			return fmt.Errorf("invalid catalog object name")
		}
	}

	return nil
}

func writeResponse(writer *bufio.Writer, success bool) error {
	if success {
		if err := writer.WriteByte(1); err != nil {
			return err
		}
	} else {
		if err := writer.WriteByte(0); err != nil {
			return err
		}
	}

	return writer.Flush()
}
