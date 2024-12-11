package peerhandler

import (
	"fmt"
	"log/slog"

	"github.com/libsv/go-p2p"
	"github.com/libsv/go-p2p/wire"
)

func StartPeerManager() error {
	pm := p2p.NewPeerManager(slog.Default(), wire.TestNet3)

	peerHandler := NewPeerHandler()

	peer1, err := p2p.NewPeer(slog.Default(), "localhost:18332", peerHandler, wire.TestNet3)
	if err != nil {
		return fmt.Errorf("error creating peer %s: %v", "1", err)
	}

	if err = pm.AddPeer(peer1); err != nil {
		return fmt.Errorf("error adding peer %s: %v", "1", err)
	}

	peer2, err := p2p.NewPeer(slog.Default(), "localhost:48332", peerHandler, wire.TestNet3)
	if err != nil {
		return fmt.Errorf("error creating peer %s: %v", "2", err)
	}

	if err = pm.AddPeer(peer2); err != nil {
		return fmt.Errorf("error adding peer %s: %v", "2", err)
	}

	return nil
}
