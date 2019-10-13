package structures

//Client is basically every user in blockchain
type Client struct {
	Name           string
	Quit           chan bool
	ListeningIP    string
	ListeningPort  string
	ConnectedPeers []Client
	AccountBalance int
}

//Transaction is basically every user to other
type Transaction struct {
	From   string
	To     chan bool
	Amount string
}
