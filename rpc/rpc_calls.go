package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kuba/nodes/utils"
	"net/http"
	"time"

	"github.com/ordishs/go-bitcoin"
)

func InvalidateAndReconsider() {
	bitcoind, err := bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		panic(err)
	}

	hashes, err := bitcoind.Generate(200)
	if err != nil {
		panic(err)
	}

	fmt.Println("200 blocks generated")

	tips, err := bitcoind.GetChainTips()
	if err != nil {
		panic(err)
	}

	utils.Print("chain tips", tips)

	// info, err := bitcoind.GetInfo()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// utils.Print("Info", info)

	blockHash := hashes[190]
	RpcCall("invalidateblock", []interface{}{blockHash})

	fmt.Println("100 blocks invalidated")
	time.Sleep(5 * time.Second)

	tips, err = bitcoind.GetChainTips()
	if err != nil {
		panic(err)
	}

	utils.Print("chain tips", tips)

	mempooltxs, err := bitcoind.GetMempoolInfo()
	utils.Print("mempool info", mempooltxs)

	// info, err = bitcoind.GetInfo()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// utils.Print("Info after invalidation", info)

	forkHash, err := bitcoind.Generate(10)
	if err != nil {
		panic(err)
	}

	fmt.Println("50 blocks generated")
	time.Sleep(5 * time.Second)

	tips, err = bitcoind.GetChainTips()
	if err != nil {
		panic(err)
	}

	utils.Print("chain tips", tips)

	RpcCall("reconsiderblock", []interface{}{blockHash})

	// _, err = bitcoind.Generate(1)
	time.Sleep(10 * time.Second)

	tips, err = bitcoind.GetChainTips()
	if err != nil {
		panic(err)
	}

	utils.Print("tips", tips)

	RpcCall("preciousblock", []interface{}{forkHash[9]})
	time.Sleep(5 * time.Second)

	_, err = bitcoind.Generate(1)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	tips, err = bitcoind.GetChainTips()
	if err != nil {
		panic(err)
	}

	utils.Print("tips", tips)
}

func RpcCall(method string, params interface{}) RpcResponse {
	c := http.Client{}

	rpcR := RpcRequest{method, params, time.Now().UnixNano(), "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)

	err := jsonEncoder.Encode(rpcR)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s://%s:%d", "http", "localhost", 18332),
		payloadBuffer,
	)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("bitcoin", "bitcoin")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var rr RpcResponse

	if resp.StatusCode != 200 {
		_ = json.Unmarshal(data, &rr)
		v, ok := rr.Err.(map[string]interface{})
		if ok {
			err = errors.New(v["message"].(string))
		} else {
			err = errors.New("HTTP error: " + resp.Status)
		}
		if err != nil {
			panic(err)
		}
	}

	err = json.Unmarshal(data, &rr)
	if err != nil {
		panic(err)
	}

	fmt.Println("successful call with method: ", method)
	return rr
}
