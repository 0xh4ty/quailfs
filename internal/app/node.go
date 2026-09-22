package app

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/0xh4ty/quailfs/internal/storage"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"os"
	"path/filepath"
)

func RunNode() {
	var bootstrapNodes []string
	var privKey crypto.PrivKey
	var dbPrivKey crypto.PrivKey
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	nodeDataDir := filepath.Join(homeDir, ".node-data")

	err = os.MkdirAll(nodeDataDir, 0700)
	if err != nil {
		panic(err)
	}

	db, err := storage.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	dbKey := uint64(0)
	dbKeyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dbKeyBytes, dbKey)

	dbPrivKeyBytes, err := storage.Get(db, "node_private_key", dbKeyBytes)
	if err != nil {
		panic(err)
	}

	if dbPrivKeyBytes == nil {
		randomness := rand.Reader
		privKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
		if err != nil {
			panic(err)
		}

		dbPrivKeyBytes, err = crypto.MarshalPrivateKey(privKey)
		if err != nil {
			panic(err)
		}

		err = storage.Put(db, "node_private_key", dbKeyBytes, dbPrivKeyBytes)
		if err != nil {
			panic(err)
		}
	} else {
		dbPrivKey, err = crypto.UnmarshalPrivateKey(dbPrivKeyBytes)
		if err != nil {
			panic(err)
		}

		privKey = dbPrivKey
	}

	node, err := libp2p.New(libp2p.Identity(privKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("PeerID: ", node.ID())

	if len(bootstrapNodes) == 0 {
		startPeer(ctx, node)
	} else {
		startPeerWithBootstrapNodes(ctx, node, bootstrapNodes)
	}

	select {}
}

func startPeer(ctx context.Context, node host.Host) {

}

func startPeerWithBootstrapNodes(ctx context.Context, node host.Host, bootstrapNodes []string) {

}
