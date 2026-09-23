package app

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/0xh4ty/quailfs/internal/config"
	"github.com/0xh4ty/quailfs/internal/storage"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	"os"
	"path/filepath"
)

func RunNode(enableRelay bool) {
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
	defer db.Close()

	dbKey := uint64(0)
	dbKeyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dbKeyBytes, dbKey)

	dbPrivKeyBytes, err := storage.Get(
		db,
		"node_private_key",
		dbKeyBytes,
	)
	if err != nil {
		panic(err)
	}

	if dbPrivKeyBytes == nil {
		randomness := rand.Reader

		privKey, _, err = crypto.GenerateKeyPairWithReader(
			crypto.RSA,
			2048,
			randomness,
		)
		if err != nil {
			panic(err)
		}

		dbPrivKeyBytes, err = crypto.MarshalPrivateKey(privKey)
		if err != nil {
			panic(err)
		}

		err = storage.Put(
			db,
			"node_private_key",
			dbKeyBytes,
			dbPrivKeyBytes,
		)
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

	libp2pOptions := []libp2p.Option{
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/4001",
			"/ip4/0.0.0.0/udp/4001/quic-v1",
		),
		libp2p.EnableAutoNATv2(),
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
	}

	if enableRelay {
		libp2pOptions = append(
			libp2pOptions,
			libp2p.EnableRelayService(),
		)
	}

	node, err := libp2p.New(libp2pOptions...)
	if err != nil {
		panic(err)
	}
	defer node.Close()

	fmt.Println("PeerID:", node.ID())

	log.Println("This node's multiaddresses:")
	for _, la := range node.Addrs() {
		log.Printf(" - %v\n", la)
	}
	log.Println()

	node.SetStreamHandler(
		"/quailfs/1.0.0",
		handleStream,
	)

	if enableRelay {
		log.Println("Relay service enabled")
	} else {
		log.Println("Relay service disabled")
	}

	bootstrapConfigPath := filepath.Join(nodeDataDir, "bootstrap.config")
	bootstrapNodes, err := config.ReadBootstrapConfig(bootstrapConfigPath)

	if len(bootstrapNodes) == 0 {
		startPeer(ctx, node)
	} else {
		startPeerWithBootstrapNodes(ctx, node, bootstrapNodes)
	}

	select {}
}

func startPeer(
	ctx context.Context,
	node host.Host,
) {
}

func startPeerWithBootstrapNodes(
	ctx context.Context,
	node host.Host,
	bootstrapNodes []string,
) {
}

func handleStream(s network.Stream) {
	defer s.Close()

	log.Println("Got a new stream!")

	reader := bufio.NewReader(s)
	writer := bufio.NewWriter(s)

	_ = reader
	_ = writer
}
