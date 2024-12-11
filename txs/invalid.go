package txs

import (
	"fmt"
	"kuba/nodes/utils"

	"github.com/ordishs/go-bitcoin"
)

func InvalidTxs() error {
	// connect to bitcoin node
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind instance: %w", err)
	}

	bitcoind2, err := bitcoin.New("localhost", 48332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind instance: %w", err)
	}

	// get info from the bitcoin node
	info, err := bitcoind.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get bitcoind info: %w", err)
	}
	utils.Print("INFO", info)

	// create a new random address
	address, err := bitcoind.GetNewAddress()
	if err != nil {
		return err
	}
	utils.Print("address", address)

	// dump xPriv of the newly created address
	privKey, err := bitcoind.DumpPrivKey(address)
	if err != nil {
		return err
	}
	utils.Print("privKey", privKey)

	// generate 200 blocks to get SATS from coinbase
	hashes, err := bitcoind.Generate(200)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	utils.Print("num of hashes generated", len(hashes))

	// send a tx
	txID, err := bitcoind.SendToAddress(address, 0.001)
	if err != nil {
		return err
	}
	utils.Print("txID", txID)

	// get rawTx by txID
	rawTx, err := bitcoind.GetRawTransaction(txID)
	if err != nil {
		return err
	}
	utils.Print("rawTx", rawTx.Hex)

	// mine the tx in the block
	hash, err := bitcoind.Generate(1)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	blockHash := hash[0]
	utils.Print("block hash", blockHash)

	// get block data
	// blockData, err := bitcoind.GetBlock(blockHash)
	// if err != nil {
	// 	return err
	// }
	// utils.Print("Block data", blockData)

	// list utxos from address
	// utxos, err := bitcoind.ListUnspent([]string{address})
	// if err != nil {
	// 	return err
	// }
	// utils.Print("UTXOS", utxos)

	rawTxString, err := bitcoind.SendRawTransactionWithoutFeeCheckOrScriptCheck("0100000001358eb38f1f910e76b33788ff9395a5d2af87721e950ebd3d60cf64bb43e77485010000006a47304402203be8a3ba74e7b770afa2addeff1bbc1eaeb0cedf6b4096c8eb7ec29f1278752602205dc1d1bedf2cab46096bb328463980679d4ce2126cdd6ed191d6224add9910884121021358f252895263cd7a85009fcc615b57393daf6f976662319f7d0c640e6189fcffffffff02bf010000000000001976a91449f066fccf8d392ff6a0a33bc766c9f3436c038a88acfc080000000000001976a914a7dcbd14f83c564e0025a57f79b0b8b591331ae288ac00000000")
	if err != nil {
		fmt.Printf("failed to send raw tx to node: %s", err.Error())
		// return fmt.Errorf("failed to send raw tx to node: %w", err)
	}
	utils.Print("rawTx", rawTxString)

	_, err = bitcoind2.SendRawTransactionWithoutFeeCheckOrScriptCheck("0100000001358eb38f1f910e76b33788ff9395a5d2af87721e950ebd3d60cf64bb43e77485010000006a47304402203be8a3ba74e7b770afa2addeff1bbc1eaeb0cedf6b4096c8eb7ec29f1278752602205dc1d1bedf2cab46096bb328463980679d4ce2126cdd6ed191d6224add9910884121021358f252895263cd7a85009fcc615b57393daf6f976662319f7d0c640e6189fcffffffff02bf010000000000001976a91449f066fccf8d392ff6a0a33bc766c9f3436c038a88acfc080000000000001976a914a7dcbd14f83c564e0025a57f79b0b8b591331ae288ac00000000")
	if err != nil {
		fmt.Printf("failed to send raw tx to node: %s", err.Error())
		// return fmt.Errorf("failed to send raw tx to node: %w", err)
	}

	// generate 100 blocks to get discarded from mempool - doesn't work
	// hashes, err = bitcoind.Generate(201)
	// if err != nil {
	// 	return fmt.Errorf("failed to generate blocks %w", err)
	// }
	// utils.Print("num of hashes generated", len(hashes))

	return nil
}
