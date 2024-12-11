package chainhash

import (
	"fmt"

	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/ordishs/go-utils"
)

type Status struct {
	CompetingTxs []string
}

func TestChainhash() {
	foundStatus := make(map[chainhash.Hash]Status, 0)
	bytea := []byte{
		0xa7, 0xfd, 0x98, 0xbd, 0x37, 0xf9, 0xb3, 0x87,
		0xdb, 0xef, 0x4f, 0x1a, 0x4e, 0x47, 0x90, 0xb9,
		0xa0, 0xd4, 0x8f, 0xb7, 0xbb, 0xb7, 0x74, 0x55,
		0xe8, 0xf3, 0x9d, 0xf0, 0xf8, 0x90, 0x9d, 0xb7,
	}

	// bytea2 := []byte{
	// 	0xa7, 0xfd, 0x98, 0xbd, 0x37, 0xf9, 0xb3, 0x87,
	// 	0xdb, 0xef, 0x4f, 0x1a, 0x4e, 0x47, 0x90, 0xb9,
	// 	0xa0, 0xd4, 0x8f, 0xb7, 0xbb, 0xb7, 0x74, 0x55,
	// 	0xe8, 0xf3, 0x9d, 0xf0, 0xf8, 0x90, 0x9d, 0xb7,
	// }

	hash, err := chainhash.NewHash(bytea)
	if err != nil {
		panic(err)
	}
	fmt.Println("hash", hash.String())

	hash2, err := chainhash.NewHash(bytea)
	if err != nil {
		panic(err)
	}
	fmt.Println("hash2", hash2.String())
	fmt.Println("both hashes equal: ", *hash == *hash2)

	status, found := foundStatus[*hash]
	if !found {
		fmt.Println("Status: ", status.CompetingTxs)
	}

	foundStatus[*hash] = Status{CompetingTxs: []string{"1234"}}
	foundStatus[*hash2] = Status{CompetingTxs: []string{"5678"}}

	status, found = foundStatus[*hash]
	if found {
		fmt.Println("Status hash 1: ", status.CompetingTxs)
	}
	status, found = foundStatus[*hash2]
	if found {
		fmt.Println("Status hash 2: ", status.CompetingTxs)
	}

	txIds := []string{"b79d90f8f09df3e85574b7bbb78fd4a0b990474e1a4fefdb87b3f937bd98fda7", "8e75ae10f86d8a43044a54c3c57d660d20cdb74e233be4b5c90ba752ebdc7e88", "d64adfce6b105dc6bdf475494925bf06802a41a0582586f33c2b16d537a0b7b6"}
	for _, s := range txIds {
		h, err := chainhash.NewHashFromStr(s)
		if err != nil {
			panic(err)
		}

		fmt.Println("Hash: ", h.String())
	}
}

func TestGoUtils() {
	bytea := []byte{
		0xa7, 0xfd, 0x98, 0xbd, 0x37, 0xf9, 0xb3, 0x87,
		0xdb, 0xef, 0x4f, 0x1a, 0x4e, 0x47, 0x90, 0xb9,
		0xa0, 0xd4, 0x8f, 0xb7, 0xbb, 0xb7, 0x74, 0x55,
		0xe8, 0xf3, 0x9d, 0xf0, 0xf8, 0x90, 0x9d, 0xb7,
	}

	hash, err := chainhash.NewHash(bytea)
	if err != nil {
		panic(err)
	}
	fmt.Println("hash", hash.String())

	fmt.Println("hash with go-utils:", utils.ReverseAndHexEncodeSlice(hash.CloneBytes()))
}
