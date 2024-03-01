package main

import (
	"fmt"

	mpt "github.com/ysfkel/merkle-patrica-trie/trie"
)

func main() {
	trie := mpt.NewTrie()
	hash0 := trie.Hash()

	trie.Put([]byte{1, 2, 3, 4}, []byte("hello"))
	hash1 := trie.Hash()

	trie.Put([]byte{1, 2, 3, 4, 5}, []byte("world"))
	hash2 := trie.Hash()

	trie.Put([]byte{1, 2, 4, 5, 6, 7, 5}, []byte("trie"))
	hash3 := trie.Hash()

	fmt.Println("", hash0, "\n", hash1, "\n", hash2, "\n", hash3)
	bytes, exists := trie.Get([]byte{1, 2, 4, 5, 6, 7, 5})
	fmt.Println(string(bytes), exists)
}
