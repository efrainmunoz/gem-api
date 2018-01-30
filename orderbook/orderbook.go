package orderbook

import "time"

// GLOBAL VARS
var orderbook = make(map[string]Orderbook)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)

// STATE
func InitState() {
	var state = make(map[string]Orderbook)

	for {
		select {
		case read := <-readsAll:
			read.resp <- state

		case read := <-readsOne:
			read.resp <- state[read.key]

		case write := <-writes:
			state[write.key] = write.val
			write.resp <- true
		}
	}
}

// WRITE new tickers
func Write(pair string, result OrderbookResponse) {

	orderbook := Orderbook{
		Asks:      result.Asks,
		Bids:      result.Bids,
		Timestamp: time.Now().Unix()}

	write := &writeOp{
		key:  pair,
		val:  orderbook,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}
