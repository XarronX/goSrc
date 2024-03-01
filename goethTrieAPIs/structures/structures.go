package structures

type SetKeyVal struct {
	Key string `json:"key"`
	Val Values `json:"val"`
}

type GetKeyVal struct {
	Key string `json:"key"`
}

type Hash struct {
	Hash string `json:"hash"`
}

type MerkleVerifcator struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

type MerkleVerifcatorResult struct {
	Hash     string `json:"hash"`
	Key      string `json:"key"`
	Verified bool   `json:"verified"`
}

type Values struct {
	Balance string `json:"balance"`
	Nonce   string `json:"nounce"`
}
