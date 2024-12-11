package beef

import (
	"fmt"
	"kuba/nodes/utils"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/ordishs/go-bitcoin"
)

func BumpGen() error {
	// connect to bitcoin node
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
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
	utils.Print("rawTx", rawTx)

	// mine the tx in the block
	hash, err := bitcoind.Generate(1)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	blockHash := hash[0]
	utils.Print("block hash", blockHash)

	// get block data
	blockData, err := bitcoind.GetBlock(blockHash)
	if err != nil {
		return err
	}
	utils.Print("Block data", blockData)

	// list utxos from address
	utxos, err := bitcoind.ListUnspent([]string{address})
	if err != nil {
		return err
	}
	utils.Print("UTXOS", utxos)

	// create BUMP
	var merkleHashes []*chainhash.Hash
	var txIndex uint64

	for i, txid := range blockData.Tx {
		if txid == txID {
			txIndex = uint64(i)
		}
		h, err := chainhash.NewHashFromStr(txid)
		if err != nil {
			return err
		}
		merkleHashes = append(merkleHashes, h)
	}

	merkleTree := bc.BuildMerkleTreeStoreChainHash(merkleHashes)

	bump, err := bc.NewBUMPFromMerkleTreeAndIndex(blockData.Height, merkleTree, txIndex)
	if err != nil {
		return err
	}
	utils.Print("BUMP", bump)

	merkleRoot, err := bump.CalculateRootGivenTxid(txID)
	if err != nil {
		return err
	}
	utils.Print("MerkleRoot from BUMP", merkleRoot)

	return nil
}
