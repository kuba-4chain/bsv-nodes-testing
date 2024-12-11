package main

import (
	"fmt"
	"kuba/nodes/beef"
	peerhandler "kuba/nodes/peer_handler"
	"kuba/nodes/txs"
	"kuba/nodes/zmq"
	"os"
	"time"
)

func main() {
	// chainhash.TestChainhash()
	// PeerManager()
	Zmq()
	DoubleSpend()
	// Beef()
	// Regular()
	// rpc.InvalidateAndReconsider()

	time.Sleep(30 * time.Second)
	// time.Sleep(5 * time.Minute)

	fmt.Println("Successful shutdown")
}

func PeerManager() {
	err := peerhandler.StartPeerManager()
	if err != nil {
		fmt.Printf("\nerror starting Peer Manager: %+v\n", err)
	}
}

func Zmq() {
	// node1
	err := zmq.RunZmq(1, 28332)
	if err != nil {
		fmt.Printf("\nerror running ZMQ %d: %+v\n", 1, err)
	}

	// node2
	// err = zmq.RunZmq(2, 38332)
	// if err != nil {
	// 	fmt.Printf("\nerror running ZMQ %d: %+v\n", 2, err)
	// }
}

func DoubleSpend() {
	err := txs.DoubleSpend()
	if err != nil {
		fmt.Printf("\nerror running double spend: %+v\n", err)
	}
}

func Regular() {
	err := txs.RegularTxs()
	if err != nil {
		fmt.Printf("\nerror running regular txs: %+v\n", err)
	}
}

func Beef() {
	err := beef.BeefGen()
	if err != nil {
		fmt.Printf("\n%+v\n", err)
		os.Exit(1)
	}
}
