package main

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func main() {
	diskdb, err := leveldb.New("test", 256, 0, "", false)
	if err != nil {
		panic(fmt.Sprintf("can't create temporary database: %v", err))
	}

	d := trie.NewDatabase(diskdb)

	t, err := trie.New(common.Hash{}, d)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	// kvs := d.DiskDB()

	// fmt.Println(kvs)

	// memdb := memorydb.New()

	// tr, err := trie.New(common.Hash{}, trie.NewDatabase(memdb))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }

	// tr.Update([]byte{1, 2, 3}, []byte("Hello"))

	// hash, i, err := tr.Commit(nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }

	// fmt.Println(hash, i)

	// db, err := ldb.OpenFile("go/src/goeth/test", &opt.Options{})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }
	// defer db.Close()

	t.Update([]byte{1, 2, 3}, []byte("Hello"))
	t.Update([]byte{1, 3, 5}, []byte("Hell"))
	t.Update([]byte{1, 4, 6}, []byte("Hellow"))
	t.Update([]byte{1, 4, 9}, []byte("He"))
	t.Update([]byte{1, 4, 8}, []byte("Hel"))
	t.Update([]byte{1, 5, 2}, []byte("Hi"))
	t.Update([]byte{1, 5, 8}, []byte("test1"))
	t.Update([]byte{2, 5, 8}, []byte("test2"))
	t.Update([]byte{3, 5, 8}, []byte("test3"))
	t.Update([]byte{4, 5, 8}, []byte("test4"))
	t.Update([]byte{5, 5, 8}, []byte("test5"))

	// fmt.Println(t.Get([]byte{5, 5, 8}))

	// hash, i, err := t.Commit(nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }
	// fmt.Println(hash, i)

	// dbval := d.Nodes()

	// fmt.Println(dbval)

	// fmt.Println(dbval[0].UnmarshalText([]byte(dbval[0].String())))

	// fmt.Println(iterateKeys(kvs.NewIterator(nil, []byte("Hi"))))

	// kvs := d.DiskDB()

	// exists, err := kvs.Has([]byte{1, 2, 3})
	// fmt.Println(exists)

	// keys := iterateKeys(kvs.NewIterator([]byte{}, []byte{1, 2, 3}))

	// fmt.Println(string(t.Get([]byte{1, 2, 3})))/ kvs := d.DiskDB()

	// exists, err := kvs.Has([]byte{1, 2, 3})
	// fmt.Println(exists)

	// keys := iterateKeys(kvs.NewIterator([]byte{}, []byte{1, 2, 3}))

	// fmt.Println(string(t.Get([]byte{1, 2, 3})))

	// db, err := ldb.Open(storage.NewMemStorage(), nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }
	// fmt.Println(db.Has([]byte{5, 5, 8}, &opt.ReadOptions{}))

	to_encode := []string{strconv.Itoa(1), strconv.Itoa(2)}

	b, err := rlp.EncodeToBytes(to_encode)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Printf("%X\n", b)

	var decoded interface{}

	err = rlp.DecodeBytes(b, &decoded)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Println(decoded)
}

// func iterateKeys(i ethdb.Iterator) []string {
// 	keys := []string{}
// 	for i.Next() {
// 		keys = append(keys, string(i.Key()))
// 	}

// 	return keys
// }
