package structures
import blockchain "github.com/syedaraiz/assignment02IBC/Blockchain"
//Client is basically every user in blockchain
type Client struct {
	Name           string
	Quit           chan bool
	ListeningIP    string
	ListeningPort  string
	ConnectedPeers []Client
	AccountBalance int
	ChainHead *blockchain.Block
}

//Transaction is basically every user to other
type Transaction struct {
	From   string
	To     chan bool
	Amount string
}
