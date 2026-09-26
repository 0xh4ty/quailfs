package network

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const PutShardProtocol = "/quailfs/put-shard/1.0.0"

func PutShard(ctx context.Context, h host.Host, peerID peer.ID, shardName string, shard []byte) error {
	stream, err := h.NewStream(ctx, peerID, PutShardProtocol)
	if err != nil {
		return fmt.Errorf("open PUT_SHARD stream: %w", err)
	}
	defer stream.Close()

	writer := bufio.NewWriter(stream)
	reader := bufio.NewReader(stream)

	if err := writeString(writer, shardName); err != nil {
		return fmt.Errorf("write shard name: %w", err)
	}

	if err := binary.Write(writer, binary.BigEndian, uint64(len(shard))); err != nil {
		return fmt.Errorf("write shard size: %w", err)
	}

	if _, err := writer.Write(shard); err != nil {
		return fmt.Errorf("write shard: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush shard: %w", err)
	}

	success, err := readResponse(reader)
	if err != nil {
		return fmt.Errorf("read PUT_SHARD response: %w", err)
	}

	if !success {
		return fmt.Errorf("peer rejected shard")
	}

	return nil
}

func writeString(writer *bufio.Writer, value string) error {
	if err := binary.Write(
		writer,
		binary.BigEndian,
		uint64(len(value)),
	); err != nil {
		return err
	}

	_, err := writer.WriteString(value)
	return err
}

func readResponse(reader *bufio.Reader) (bool, error) {
	status, err := reader.ReadByte()
	if err != nil {
		return false, err
	}

	switch status {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid response status: %d", status)
	}
}
