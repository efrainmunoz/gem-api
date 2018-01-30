package ticker

// MODELS
type Ticker struct {
	LastPrice string `json:"lastprice"`
	BestBid   string `json:"bestbid"`
	BestAsk   string `json:"bestask"`
	Timestamp int64  `json:"timestamp"`
}

type TickerResponse struct {
	Bid  string `json:"bid"`
	Ask  string `json:"ask"`
	Last string `json:"last"`
}

type readAllOp struct {
	resp chan map[string]Ticker
}

type readOneOp struct {
	key  string
	resp chan Ticker
}

type writeOp struct {
	key  string
	val  Ticker
	resp chan bool
}

