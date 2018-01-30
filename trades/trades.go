package trades

// GLOBAL VARS
var trades = make(map[string]Trades)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)

// STATE
func InitState() {
	var state = make(map[string]Trade)

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
func Write(sg3Pair string, trade Trade) {

	write := &writeOp{
		key:  sg3Pair,
		val: trade,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}
