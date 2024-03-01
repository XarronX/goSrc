package keyValDB

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
)

func GetKeyValDB() *leveldb.Database {
	diskdb, err := leveldb.New("test", 256, 0, "", false)
	if err != nil {
		panic(fmt.Sprintf("can't create temporary database: %v", err))
	}

	return diskdb
}
