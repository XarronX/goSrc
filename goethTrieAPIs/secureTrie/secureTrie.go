package secureTrie

import (
	"fmt"

	tree "github.com/ashutosh/goethTrieAPIs/trie"
	"github.com/ethereum/go-ethereum/trie"
)

type TreeImplementation interface {
	Init()
	Set([]byte, []byte)
	Get([]byte) []byte
	GetTree() *SecureTree
	// GetHash() string
}

type SecureTree struct {
	secureTrie *trie.SecureTrie
}

var st *SecureTree

func Init() {
	s, err := trie.NewSecure(tree.GetTree().GetHash(), tree.GetDB())
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	st = &SecureTree{secureTrie: s}
}

func (t *SecureTree) Set(key, val []byte) {
	t.secureTrie.Update(key, val)
	h, i, err := t.secureTrie.Commit(nil)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(h, i)
}

func (t *SecureTree) GetValue(key []byte) []byte {
	return t.secureTrie.Get(key)
}

func GetSecureTree() *SecureTree {
	return st
}
