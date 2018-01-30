package trades

// MODELS
type Trade struct {
	Price   string `json:"price"`
	Volume  string `json:"volume"`
	TradeAction string `json:"trade-action"`
}

type Trades struct {
	Trades    []Trade `json:"trades"`
	Timestamp int64   `json:"timestamp"`
}

type readAllOp struct {
	resp chan map[string]Trade
}

type readOneOp struct {
	key  string
	resp chan Trade
}

type writeOp struct {
	key  string
	val  Trade
	resp chan bool
}

// Models to deal with the third-party api responses
type trade struct {
	ID        int     `json:"Id"`
	TimeStamp string  `json:"TimeStamp"`
	Quantity  float64 `json:"Quantity"`
	Price     float64 `json:"Price"`
	Total     float64 `json:"Total"`
	FillType  string  `json:"FillType"`
	OrderType string  `json:"OrderType"`
}

type TradesResult []trade

type TradesResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Result  TradesResult `json:"result"`
}
