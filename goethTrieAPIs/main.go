package main

import (
	"github.com/ashutosh/goethTrieAPIs/server"
	"github.com/ashutosh/goethTrieAPIs/trie"
)

func main() {
	trie.Init()
	server.Server()
}
