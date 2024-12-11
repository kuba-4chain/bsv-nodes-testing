package zmq

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"kuba/nodes/utils"

	"github.com/ordishs/go-bitcoin"
)

const (
	hashtxTopic                    = "hashtx2"
	invalidtxTopic                 = "invalidtx"
	discardedFromMempoolTopic      = "discardedfrommempool"
	discardedFromBlockMempoolTopic = "removedfrommempoolblock"
)

func RunZmq(number, port int) error {
	zmq := bitcoin.NewZMQ("localhost", port)

	ch := make(chan []string)

	go reader(ch, number)

	if err := zmq.Subscribe(hashtxTopic, ch); err != nil {
		return err
	}

	if err := zmq.Subscribe(invalidtxTopic, ch); err != nil {
		return err
	}

	if err := zmq.Subscribe(discardedFromMempoolTopic, ch); err != nil {
		return err
	}

	if err := zmq.Subscribe(discardedFromBlockMempoolTopic, ch); err != nil {
		return err
	}

	return nil
}

func reader(ch chan []string, number int) {
	name := fmt.Sprintf("ZMQ %d\t", number)

	for c := range ch {
		utils.Print(name+"msg RAW", c[1])
		switch c[0] {
		case "hashtx2":
			utils.Print(name+"ACCEPTED_BY_NETWORK", c[1])

		case invalidtxTopic:
			// c[1] is lots of info about the tx in json format encoded in hex
			jsonHex, err := hex.DecodeString(c[1])
			if err != nil {
				fmt.Printf("ZMQ %d: Error reading invalid tx message\n", number)
				continue
			}
			utils.Print(name+"invalid tx", string(jsonHex))

			var txInfo *ZMQTxInfo
			txInfo, err = parseTxInfo(c)
			if err != nil {
				fmt.Printf("ZMQ %d: Error parsing invalid tx info\n", number)
				continue
			}
			utils.Print(name+"invalid info", txInfo)

		case discardedFromMempoolTopic:
			txDiscardedInfo, err := hex.DecodeString(c[1])
			if err != nil {
				fmt.Printf("ZMQ %d: Error parsing discarded from mempool\n", number)
				continue
			}
			utils.Print(name+"discarded from mempool", string(txDiscardedInfo))

		case discardedFromBlockMempoolTopic:
			txDiscardedInfo, err := hex.DecodeString(c[1])
			if err != nil {
				fmt.Printf("ZMQ %d: Error parsing discarded from mempool\n", number)
				continue
			}
			utils.Print(name+"discarded from BLOCK mempool", string(txDiscardedInfo))

		default:
			utils.Print(name+"unhandled zmq msg", c[1])
		}
	}
	fmt.Printf("ZQM %d channel closed\n", number)
}

func parseTxInfo(c []string) (*ZMQTxInfo, error) {
	var txInfo ZMQTxInfo
	txInfoBytes, err := hex.DecodeString(c[1])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(txInfoBytes, &txInfo)
	if err != nil {
		return nil, err
	}
	return &txInfo, nil
}
