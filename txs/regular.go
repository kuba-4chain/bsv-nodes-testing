package txs

import (
	"fmt"
	"kuba/nodes/utils"

	"github.com/ordishs/go-bitcoin"
)

func RegularTxs() error {
	// connect to bitcoin node
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind instance: %w", err)
	}

	// generate 200 blocks to get SATS from coinbase
	hashes, err := bitcoind.Generate(200)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	utils.Print("num of hashes generated", len(hashes))

	// create a new random address
	address, err := bitcoind.GetNewAddress()
	if err != nil {
		return err
	}
	utils.Print("address", address)

	// send a tx
	txID, err := bitcoind.SendToAddress(address, 0.001)
	if err != nil {
		return err
	}
	utils.Print("txID", txID)

	return nil
}
