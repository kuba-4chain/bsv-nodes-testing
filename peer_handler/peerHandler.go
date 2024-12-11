package peerhandler

import (
	"context"
	"fmt"

	"github.com/libsv/go-p2p"
	"github.com/libsv/go-p2p/wire"
)

type PeerHandler struct {
	cancelAll context.CancelFunc
	ctx       context.Context
}

func NewPeerHandler() *PeerHandler {
	ph := &PeerHandler{}

	ctx, cancelAll := context.WithCancel(context.Background())
	ph.cancelAll = cancelAll
	ph.ctx = ctx

	return ph
}

// HandleTransactionSent is called when a transaction is sent to a peer.
func (m *PeerHandler) HandleTransactionSent(msg *wire.MsgTx, peer p2p.PeerI) error {
	hash := msg.TxHash().String()
	fmt.Printf("PEER_HANDLER: tx: %s SENT_TO_NETWORK to peer: %s", hash, peer.String())
	return nil
}

// HandleTransactionAnnouncement is a message sent to the PeerHandler when a transaction INV message is received from a peer.
func (m *PeerHandler) HandleTransactionAnnouncement(msg *wire.InvVect, peer p2p.PeerI) error {
	fmt.Printf("PEER_HANDLER: received INV for tx: %s from peer: %s - SEEN_ON_NETWORK", msg.Hash.String(), peer.String())
	return nil
}

// HandleTransactionRejection is called when a transaction is rejected by a peer.
func (m *PeerHandler) HandleTransactionRejection(rejMsg *wire.MsgReject, peer p2p.PeerI) error {
	fmt.Printf("PEER_HANDLER: rejected tx: %s, by peer: %s, reason: %s - REJECTED", rejMsg.Hash.String(), peer.String(), rejMsg.Reason)
	return nil
}

// HandleTransactionsGet is called when a peer requests a transaction.
func (m *PeerHandler) HandleTransactionsGet(msgs []*wire.InvVect, peer p2p.PeerI) ([][]byte, error) {
	fmt.Printf("PEER_HANDLER: requested number of txs: %d, from peer: %s - REQUEST_BY_NETWORK, sending nothing", len(msgs), peer.String())
	return [][]byte{}, nil
}

// HandleTransaction is called when a transaction is received from a peer.
func (m *PeerHandler) HandleTransaction(msg *wire.MsgTx, peer p2p.PeerI) error {
	hash := msg.TxHash().String()
	fmt.Printf("PEER_HANDLER: got tx: %s, from peer: %s - SEEN_ON_NETWORK", hash, peer.String())
	return nil
}

// HandleBlockAnnouncement is called when a block INV message is received from a peer.
func (m *PeerHandler) HandleBlockAnnouncement(_ *wire.InvVect, _ p2p.PeerI) error {
	return nil
}

// HandleBlock is called when a block is received from a peer.
func (m *PeerHandler) HandleBlock(_ wire.Message, _ p2p.PeerI) error {
	return nil
}

func (m *PeerHandler) Shutdown() {
	if m.cancelAll != nil {
		m.cancelAll()
	}
}
