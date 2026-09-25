package place

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"math/rand"
)

func Place(livePeers []peer.ID, objectNames []string) map[peer.ID][]string {
	placements := make(map[peer.ID][]string)

	if len(livePeers) == 0 || len(objectNames) == 0 {
		return placements
	}

	for _, peerID := range livePeers {
		placements[peerID] = []string{}
	}

	peers := make([]peer.ID, len(livePeers))
	copy(peers, livePeers)

	rand.Shuffle(
		len(peers),
		func(i, j int) {
			peers[i], peers[j] = peers[j], peers[i]
		},
	)

	for i, objectName := range objectNames {
		peerID := peers[i%len(peers)]
		placements[peerID] = append(
			placements[peerID],
			objectName,
		)
	}

	return placements
}
