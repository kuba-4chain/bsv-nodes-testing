package txs

import (
	"errors"
	"fmt"
	"kuba/nodes/rpc"
	"kuba/nodes/utils"
	"time"

	"github.com/ordishs/go-bitcoin"
)

func DoubleSpend() error {
	// connect to bitcoin node
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind instance: %w", err)
	}

	// connnect to a second node for double spend
	bitcoind2, err := bitcoin.New("localhost", 48332, "bitcoin", "bitcoin", false)
	if err != nil {
		return fmt.Errorf("failed to create bitcoind2 instance: %w", err)
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

	time.Sleep(2 * time.Second)
	fmt.Println("================ 200 blocks generated =================")
	fmt.Println()

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

	time.Sleep(2 * time.Second)
	fmt.Println("================ transaction sent to address and mined =================")
	fmt.Println()

	// list utxos from address
	utxos, err := bitcoind.ListUnspent([]string{address})
	if err != nil {
		return err
	}
	if len(utxos) == 0 {
		return errors.New("no UTXO found")
	}
	// utils.Print("UTXOS", utxos)

	// first tx using the same input
	firstAddr, err := bitcoind.GetNewAddress()
	if err != nil {
		return err
	}
	// utils.Print("First Address", firstAddr)

	tx, err := utils.CreateTx(privKey, firstAddr, utxos[0])
	if err != nil {
		return err
	}
	utils.Print("first tx hex", tx.String())

	// second tx using the same input
	secondAddr := "mnGP8pC9JcjxwSi2zjgBnZCnvX1odxqfLz"
	// utils.Print("Second Address", secondAddr)

	tx2, err := utils.CreateTx(privKey, secondAddr, utxos[0])
	if err != nil {
		return err
	}
	utils.Print("second tx hex", tx2.String())

	time.Sleep(2 * time.Second)
	fmt.Println("================ double spend hexes prepared =================")
	fmt.Println()

	// post 1st tx to node
	txid, err := bitcoind.SendRawTransaction(tx.String())
	if err != nil {
		return err
	}
	utils.Print("first txid submitted", txid)

	time.Sleep(2 * time.Second)

	// post 2nd tx to a second node
	txid2, err := bitcoind2.SendRawTransaction(tx2.String())
	if err != nil {
		fmt.Printf("error submitting 2nd tx: %v", err)
		// return err
	}
	utils.Print("second txid submitted", txid2)

	time.Sleep(2 * time.Second)
	fmt.Println("================ double spend submitted =================")
	fmt.Println()

	hashes, err = bitcoind.Generate(10)
	if err != nil {
		return fmt.Errorf("failed to generate blocks %w", err)
	}
	utils.Print("num of hashes generated", len(hashes))

	// send double spending transaction when first tx was mined
	fmt.Println("================ double spend submitted after mined =================")
	fmt.Println()

	thirdAddress := "18VWHjMt4ixHddPPbs6righWTs3Sg2QNcn"
	// utils.Print("Second Address", secondAddr)

	tx3, err := utils.CreateTx(privKey, thirdAddress, utxos[0])
	if err != nil {
		return err
	}
	utils.Print("third tx hex", tx3.String())

	// post 3rd tx to node
	txid3, err := bitcoind.SendRawTransaction(tx3.String())
	if err != nil {
		fmt.Println("error sending third tx ", err)
	}
	utils.Print("third txid submitted", txid3)

	mempool := rpc.RpcCall("getorphaninfo", nil)
	var data []byte
	err = mempool.Result.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	utils.Print("mempool1", data)

	return nil
}
