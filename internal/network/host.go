package network

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

func NewHost(privKey crypto.PrivKey, enableRelay bool) (host.Host, error) {
	options := []libp2p.Option{
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
		options = append(
			options,
			libp2p.EnableRelayService(),
		)
	}

	return libp2p.New(options...)
}
