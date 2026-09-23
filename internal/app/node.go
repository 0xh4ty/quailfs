package app

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/0xh4ty/quailfs/internal/config"
	"github.com/0xh4ty/quailfs/internal/network"
	"github.com/0xh4ty/quailfs/internal/storage"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"path/filepath"
	"time"
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

	node, err := network.NewHost(privKey, enableRelay)
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

	var kad *dht.IpfsDHT

	if len(bootstrapNodes) == 0 {
		kad, err = startPeer(node)
	} else {
		kad, err = startPeerWithBootstrapNodes(
			ctx,
			node,
			bootstrapNodes,
		)
	}
	if err != nil {
		panic(err)
	}

	defer kad.Close()

	select {}
}

func startPeer(node host.Host) (*dht.IpfsDHT, error) {

	log.Println("Starting DHT...")

	kad, err := dht.New(
		node,
		dht.Mode(dht.ModeServer),
	)
	if err != nil {
		return nil, err
	}

	log.Println("DHT started")

	return kad, nil
}

func startPeerWithBootstrapNodes(
	ctx context.Context,
	node host.Host,
	bootstrapNodes []string,
) (*dht.IpfsDHT, error) {

	var lastErr error

	kad, err := dht.New(
		node,
		dht.Mode(dht.ModeServer),
	)
	if err != nil {
		return nil, err
	}

	log.Println("DHT started")

	for _, address := range bootstrapNodes {
		log.Printf("Connecting to bootstrap node: %s\n", address)

		maddr, err := ma.NewMultiaddr(address)
		if err != nil {
			log.Printf(
				"Invalid bootstrap multiaddress %q: %v\n",
				address,
				err,
			)
			lastErr = err
			continue
		}

		addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Printf(
				"Invalid bootstrap peer address %q: %v\n",
				address,
				err,
			)
			lastErr = err
			continue
		}

		connectCtx, cancel := context.WithTimeout(
			ctx,
			10*time.Second,
		)

		err = node.Connect(connectCtx, *addrInfo)

		cancel()

		if err != nil {
			log.Printf(
				"Failed to connect to bootstrap node %s: %v\n",
				address,
				err,
			)
			lastErr = err
			continue
		}

		log.Printf(
			"Connected to bootstrap peer: %s\n",
			addrInfo.ID,
		)

		if err := kad.Bootstrap(ctx); err != nil {
			kad.Close()
			return nil, err
		}

		log.Println("DHT bootstrap completed")

		if err := verifyDHTDiscovery(
			ctx,
			kad,
			addrInfo.ID,
		); err != nil {
			kad.Close()
			return nil, err
		}

		return kad, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("no valid bootstrap nodes")
	}

	return nil, lastErr
}

func handleStream(s net.Stream) {
	defer s.Close()

	log.Println("Got a new stream!")

	reader := bufio.NewReader(s)
	writer := bufio.NewWriter(s)

	_ = reader
	_ = writer
}

func verifyDHTDiscovery(
	ctx context.Context,
	kad *dht.IpfsDHT,
	target peer.ID,
) error {

	log.Printf(
		"Verifying DHT discovery of peer: %s\n",
		target,
	)

	lookupCtx, cancel := context.WithTimeout(
		ctx,
		15*time.Second,
	)
	defer cancel()

	peers, err := kad.GetClosestPeers(
		lookupCtx,
		string(target),
	)
	if err != nil {
		return fmt.Errorf(
			"DHT peer lookup failed: %w",
			err,
		)
	}

	for _, p := range peers {
		log.Printf(
			"DHT discovered peer: %s\n",
			p,
		)

		if p == target {
			log.Printf(
				"DHT discovery verified: %s\n",
				target,
			)
			return nil
		}
	}

	return fmt.Errorf(
		"DHT lookup completed but target peer %s was not found",
		target,
	)
}
