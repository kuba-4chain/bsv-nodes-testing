package beef

import (
	"errors"
	"fmt"
	"kuba/nodes/utils"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/ordishs/go-bitcoin"
)

func BeefGen() error {
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind instance: %w", err)
	}

	address, err := bitcoind.GetNewAddress()
	if err != nil {
		return err
	}
	utils.Print("address", address)

	privKey, err := bitcoind.DumpPrivKey(address)
	if err != nil {
		return err
	}
	utils.Print("privKey", privKey)

	// info, err := bitcoind.GetInfo()
	// if err != nil {
	// 	return fmt.Errorf("failed to get bitcoind info: %w", err)
	// }
	// utils.Print("INFO", info)

	hashes, err := bitcoind.Generate(200)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	utils.Print("Hashes", len(hashes))
	return nil

	txID, err := bitcoind.SendToAddress(address, 0.001)
	if err != nil {
		return err
	}
	utils.Print("txID", txID)

	hashes, err = bitcoind.Generate(1)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	blockHash := hashes[0]
	utils.Print("block hash", blockHash)

	rawTx, err := bitcoind.GetRawTransaction(txID)
	if err != nil {
		return err
	}
	utils.Print("rawTx block hash", rawTx.BlockHash)

	if blockHash != rawTx.BlockHash {
		return errors.New("block hash mismatch")
	}

	blockData, err := bitcoind.GetBlock(blockHash)
	if err != nil {
		return err
	}
	utils.Print("block data tx", blockData.Tx)
	utils.Print("block data merkleRoot", blockData.MerkleRoot)

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
	merkleRoot, err := bump.CalculateRootGivenTxid(txID)
	if err != nil {
		return err
	}
	utils.Print("MerkleRoot", merkleRoot)

	if merkleRoot != blockData.MerkleRoot {
		return errors.New("merkle root mismatch")
	}

	utxos, err := bitcoind.ListUnspent([]string{address})
	if err != nil {
		return err
	}
	if len(utxos) == 0 {
		return errors.New("no UTXO found")
	}

	newAddress, err := bitcoind.GetNewAddress()
	if err != nil {
		return err
	}
	utils.Print("New Address", newAddress)

	tx, err := utils.CreateTx(privKey, newAddress, utxos[0])
	if err != nil {
		return err
	}
	utils.Print("tx hex", tx.String())
	utils.Print("raw tx hex", rawTx.Hex)

	beef, err := buildBeef(rawTx.Hex, bump, tx)
	if err != nil {
		return err
	}
	utils.Print("BEEF", beef)

	return nil
}

func buildBeef(rawTxInput string, bump *bc.BUMP, newTx *bt.Tx) (string, error) {
	version := "0100BEEF"
	nBumps := "01"
	bumpData, err := bump.String()
	if err != nil {
		return "", fmt.Errorf("error converting bump to string: %w", err)
	}
	nTransactions := "02"
	rawTx := rawTxInput
	hasBump := "01"
	bumpIndex := "00"
	rawTxCurrent := newTx.String()
	hasBumpCurrent := "00"

	beef := version +
		nBumps +
		bumpData +
		nTransactions +
		rawTx +
		hasBump +
		bumpIndex +
		rawTxCurrent +
		hasBumpCurrent

	return beef, nil
}
