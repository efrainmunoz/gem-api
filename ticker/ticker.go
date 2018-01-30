package ticker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GLOBAL VARS
var tickers = make(map[string]Ticker)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)

// Get a ticker from Gemini api
func GetTicker(pair string) (aTickerResponse TickerResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	url := fmt.Sprintf("https://api.gemini.com/v1/pubticker/%s", pair)

	// try to get kraken ticker
	resp, err := httpCLI.Get(url)
	if err != nil {
		return TickerResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TickerResponse{}, err
	}

	// Unmarshal the json
	tickerResponse := TickerResponse{}
	err = json.Unmarshal(body, &tickerResponse)
	if err != nil {
		return TickerResponse{}, err
	}

	return tickerResponse, nil
}


// STATE
func InitState() {
	var state = make(map[string]Ticker)

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
func Write(pair string, result TickerResponse) {
	ticker := Ticker{
		LastPrice: result.Last,
		BestBid:   result.Bid,
		BestAsk:   result.Ask,
		Timestamp: time.Now().Unix(),
	}

	write := &writeOp{
		key:  pair,
		val:  ticker,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}
