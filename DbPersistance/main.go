package main

import (
	"fmt"
	"io"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	// goldb "github.com/syndtr/goleveldb/leveldb"
)

type st struct {
	Val1 string
	Val2 string
}

func main() {
	// diskdb, err := leveldb.New("test", 256, 0, "", false)
	// if err != nil {
	// 	fmt.Printf("can't create temporary database: %v", err)
	// }

	// fmt.Println(diskdb.Path())

	// diskdb.Put([]byte("testKey"), []byte("testVal"))

	// it := diskdb.NewIterator(nil, nil)

	// for it.Next() {
	// 	fmt.Println(string(it.Key()), string(it.Value()))
	// }

	// existingDb, err := goldb.OpenFile("test", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// itr := existingDb.NewIterator(nil, nil)

	// for itr.Next() {
	// 	fmt.Println(string(itr.Key()), string(itr.Value()))
	// }

	db, err := leveldb.New("test", 256, 0, "", false)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	trieDB := trie.NewDatabase(db)
	// hash := common.Hash{}
	hash := common.HexToHash("0x993534c9ec1f458f583eba002da6e53ca25a9214862a715e57eff7944ea16ccc")

	t, err := trie.New(common.Hash{}, hash, trieDB)
	if err != nil {
		fmt.Println(err)
	}

	val := Decode(t.Get([]byte("testKey1")))

	decodedVals := val.([]interface{})
	fmt.Println(decodedVals)

	sv := &st{}
	mp := make(map[string]interface{})

	for i, v := range decodedVals {
		mp[reflect.Indirect(reflect.ValueOf(*sv)).Type().Field(i).Name] = string(v.([]byte))
	}

	// fmt.Printf("%+v\n", mp)

	// st1 := &st{Val1: "testVal110", Val2: "testVal120"}
	// st2 := &st{Val1: "testVal210", Val2: "testVal220"}
	// st3 := &st{Val1: "testVal310", Val2: "testVal320"}
	// st4 := &st{Val1: "testVal410", Val2: "testVal420"}
	// st5 := &st{Val1: "testVal510", Val2: "testVal520"}
	// st6 := &st{Val1: "testVal610", Val2: "testVal620"}
	// st7 := &st{Val1: "testVal710", Val2: "testVal720"}
	// st8 := &st{Val1: "testVal810", Val2: "testVal820"}

	// t.Update([]byte("testKey1"), Encode(st1))
	// t.Update([]byte("testKey2"), Encode(st2))
	// t.Update([]byte("testKey3"), Encode(st3))
	// t.Update([]byte("testKey4"), Encode(st4))
	// t.Update([]byte("testKey5"), Encode(st5))
	// t.Update([]byte("testKey6"), Encode(st6))
	// t.Update([]byte("testKey7"), Encode(st7))
	// t.Update([]byte("testKey8"), Encode(st8))

	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey1"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey2"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey3"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey4"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey5"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey6"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey7"))))
	// fmt.Printf("%+v\n", Decode(t.Get([]byte("testKey8"))))

	// fmt.Println(t.Commit(nil))
	// fmt.Println(trieDB.Commit(t.Hash(), true, nil))

	it := t.NodeIterator(common.Hex2Bytes(hash.Hex()))
	for i := 0; it.Next(false); i++ {
		fmt.Println(Try((it.NodeBlob())))
	}

	// sc := state.NewStateSync(hash, db, nil)

	// fmt.Println(t.Hash())

	// sc := trie.NewSync(hash, db, nil)

	// nodes, paths, codes := sc.Missing(0)
	// fmt.Println(nodes, string(paths[0][0]), codes)

	// data := trie.NewDatabase(db).Nodes()
	// fmt.Println(data)

	// val, err := db.Get([]byte("testKey"))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(val))
}

func Encode(val interface{}) []byte {
	b, err := rlp.EncodeToBytes(val)
	if err != nil {
		panic(err)
	}

	return b
}

func Decode(b []byte) interface{} {
	var decoded interface{}
	err := rlp.DecodeBytes(b, &decoded)
	if err != nil {
		panic(err)
	}

	return decoded
}

func Try(intr interface{}) (interface{}, error) {
	var reader io.Reader
	err := rlp.Decode(reader, &intr)
	if err != nil {
		return nil, err
	}

	return intr, nil
}
