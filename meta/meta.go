package meta

type Block struct {
	Index     int    `db:bIndex`
	Timestamp string `db:Timestamp`
	PrevHash  string `db:Prevhash`
	Merkle    string `db:Merkle`
	Signature string `db:Signature`
	Hash      string `db:Hash`
	TX        []Transaction
}

type BlockHeader struct {
	Index     int    `db:bIndex`
	Timestamp string `db:Timestamp`
	BPM       int    `db:BPM`
	Hash      string `db:Hash`
	PrevHash  string `db:Prevhash`
	Merkle    string
}

type Transaction struct {
	From      string
	To        string
	Timestamp string
	Signature string
	Message   string
}

type Node struct {
	IsLeader bool
	Addr     string
	Port     string
}

type TransHash struct {
	BlockHash  string
	TransHashs []string
}

type Vote struct {
	Sender string
	Hash   string
	Vote   bool
}
type ReVote struct {
	Sender string
	Vote   []Vote
	Hash   string
	V      bool
}
