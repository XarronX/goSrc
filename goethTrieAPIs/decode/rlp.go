package decode

import (
	"github.com/ashutosh/goethTrieAPIs/structures"
	"github.com/ethereum/go-ethereum/rlp"
)

func Decode(b []byte) structures.Values {
	var decoded structures.Values

	err := rlp.DecodeBytes(b, &decoded)
	if err != nil {
		panic(err)
	}

	return decoded
}
