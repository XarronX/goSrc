package api

import (
	"encoding/json"
	"net/http"

	"github.com/ashutosh/goethTrieAPIs/decode"
	"github.com/ashutosh/goethTrieAPIs/encode"
	"github.com/ashutosh/goethTrieAPIs/secureTrie"
	"github.com/ashutosh/goethTrieAPIs/structures"
	"github.com/ashutosh/goethTrieAPIs/trie"
)

func SetValue(w http.ResponseWriter, r *http.Request) {
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
}

func GetHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := trie.GetTree()
	payload := structures.Hash{
		Hash: t.GetHashStr(),
	}
	json.NewEncoder(w).Encode(payload)
	w.WriteHeader(http.StatusOK)
}

func MerkleVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mv := structures.MerkleVerifcator{}
	if err := json.NewDecoder(r.Body).Decode(&mv); err != nil {
		panic(err)
	}

	var response structures.MerkleVerifcatorResult
	if trie.GetTree().MerkleVerification(mv.Hash, mv.Key) {
		response = structures.MerkleVerifcatorResult{
			Hash:     mv.Hash,
			Key:      mv.Key,
			Verified: true,
		}
	} else {
		response = structures.MerkleVerifcatorResult{
			Hash:     mv.Hash,
			Key:      mv.Key,
			Verified: false,
		}
	}
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}

func SecureTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	secureTrie.Init()

	w.WriteHeader(http.StatusOK)
}

func GetSecureTreeVal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	kv := structures.GetKeyVal{}
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		panic(err)
	}

	st := secureTrie.GetSecureTree()

	val := st.GetValue([]byte(kv.Key))

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
}
