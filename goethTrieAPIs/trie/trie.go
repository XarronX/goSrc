package trie

import (
	"fmt"

	"github.com/ashutosh/goethTrieAPIs/keyValDB"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/trie"
)

type TreeImplementation interface {
	Set([]byte, []byte)
	Get([]byte) []byte
	GetTree() *Tree
	GetHashStr() string
	MerkleVerification(hashstr string, key []byte) bool
}

type Tree struct {
	trie *trie.Trie
}

var t *Tree
var db *trie.Database
var kvdb *leveldb.Database

func Init() {
	kvdb := keyValDB.GetKeyValDB()

	db = trie.SetValue(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	
		kv := structures.SetKeyVal{}
		t := trie.GetTree()
	
		if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
			panic(err)
		}
		t.Set([]byte(kv.Key), []byte(encode.Encode(kv.Val)))
	
		w.WriteHeader(http.StatusOK)
	}
	
	func GetValue(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		t := trie.GetTree()
	
		kv := structures.GetKeyVal{}
		if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
			panic(err)
		}
	
		val := t.Get([]byte(kv.Key))
	
		decodedVals := decode.Decode(val)
	
		payload := structures.SetKeyVal{
			Key: kv.Key,
			Val: structures.Values{
				Balance: decodedVals.Balance,
				Nonce:   decodedVals.Nonce,
			},
		}
		json.NewEncoder(w).Encode(payload)
	
		w.WriteHeader(http.StatusOK)
	}se(kvdb)

	tr, err := trie.New(common.Hash{}, db)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	t = &Tree{trie: tr}

	h, i, err := t.trie.Commit(nil)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(h, i)
}

func (tree *Tree) Set(key, val []byte) {
	tree.trie.Update(key, val)
	h, i, err := tree.trie.Commit(nil)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(h, i)
}

func (tree *Tree) Get(key []byte) []byte {
	return tree.trie.Get(key)
}

func (tree *Tree) GetHashStr() string {
	return tree.trie.Hash().String()
}

func (tree *Tree) GetHash() common.Hash {
	return tree.trie.Hash()
}

func (tree *Tree) MerkleVerification(hashstr, key string) bool {
	fmt.Println(tree.trie.Hash().String(), "\n", t.trie.Hash().String())
	hash, err := trie.VerifyProof(tree.trie.Hash(), []byte(key), kvdb)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if hashstr != string(hash) {
		return false
	}
	return true
}

func GetTree() *Tree {
	return t
}

func GetDB() *trie.Database {
	return db
}
