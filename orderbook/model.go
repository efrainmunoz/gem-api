package orderbook

// MODELS
type Order struct {
	Price    string `json:"price"`
	Volume   string `json:"volume"`
}

type Orderbook struct {
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
	Timestamp int64   `json:"timestamp"`
}

type readAllOp struct {
	resp chan map[string]Orderbook
}

type readOneOp struct {
	key  string
	resp chan Orderbook
}

type writeOp struct {
	key  string
	val  Orderbook
	resp chan bool
}

// Models to deal with the third-party api responses
type order struct {
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	Timestamp string `json:"timestamp"`
}

type OrderbookResponse struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}
