package encode

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func Encode(val interface{}) []byte {
	b, err := rlp.EncodeToBytes(val)
	if err != nil {
		panic(err)
	}

	return b
}
