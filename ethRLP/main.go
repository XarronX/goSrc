package main

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type Child struct {
	SingleStr string
	MultiStrs string
}

type Parent struct {
	Key string
	Val Child
}

func main() {
	// fmt.Println(common.HexToAddress((common.Hash{}).Hex()))
	// fmt.Println(string(crypto.Keccak256(nil)))
	// mp := map[string]interface{}{
	// 	"key": "testkey",
	// 	"val": map[string]interface{}{
	// 		"currentBlockHash":  "0X12345678defab2345433352636easad3e3e9",
	// 		"transactionHashes": []string{"0X12345678defab2345433352636easad3e3e9", "0X12345678defab2345433352636easad3e3e9"},
	// 	},
	// }
	// bytes, err := rlp.EncodeToBytes(mp["val"].(map[string]interface{})["transactionHashes"].([]string))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// ch := Child{
	// 	SingleStr: mp["val"].(map[string]interface{})["currentBlockHash"].(string),
	// 	MultiStrs: string(bytes),
	// }

	// // m := map[string]interface{}{
	// // 	"currentBlockHash": mp["val"].(map[string]interface{})["currentBlockHash"].(string),
	// // 	"transactionHashes": bytes,
	// // }

	// bytes, err = rlp.EncodeToBytes(ch)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// var intr interface{}
	// err = rlp.DecodeBytes(bytes, &intr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(intr.([]interface{})[0].([]byte)))

	// ///////////////////////////////////DECODE ENCODED HASHES ARRAY//////////////////////////////
	// var p interface{}
	// err = rlp.DecodeBytes(intr.([]interface{})[1].([]byte), &p)
	// fmt.Println(string(p.([]interface{})[0].([]byte)), string(p.([]interface{})[1].([]byte)))
	// ////////////////////////////////////////////////////////////////////////////////////////////

	fmt.Println(len(common.HexToAddress("0x2981562AD90487e8FCa217Fa51996039fff54b0f").Bytes()))

	diskdb, err := leveldb.New("test", 256, 0, "", false)
	if err != nil {
		fmt.Printf("can't create temporary database: %v", err)
	}

	storageDb := trie.NewDatabase(diskdb)

	root := common.Hash{}

	tree, err := trie.New(common.BytesToHash([]byte("Ashutosh")), root, storageDb)
	if err != nil {
		fmt.Println("Could not Create/Get trie")
		return
	}

	encodedBytes, err := rlp.EncodeToBytes(Child{
		SingleStr: "SingleStr",
		MultiStrs: "MultiStrs",
	})
	if err != nil {
		fmt.Println("Could not encode in RLP")
		return
	}

	tree.Update([]byte("TestKey"), encodedBytes)

	var child Child

	rlp.Decode(bytes.NewReader(tree.Get([]byte("TestKey"))), &child)

	fmt.Println(child)

	encodedBytes, err = rlp.EncodeToBytes(Child{
		SingleStr: "newSingleStr",
		MultiStrs: "newMultiStrs",
	})

	tree.Update([]byte("TestKey"), encodedBytes)

	rlp.Decode(bytes.NewReader(tree.Get([]byte("TestKey"))), &child)

	fmt.Println(child)
}
