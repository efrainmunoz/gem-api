package main

import (
	"github.com/efrainmunoz/gem-api/orderbook"
	"github.com/efrainmunoz/gem-api/ticker"
	"github.com/efrainmunoz/gem-api/trades"
	"github.com/efrainmunoz/gem-api/wsclient"
	"github.com/gorilla/mux"
	"net/http"
)
// MAIN
func main() {

	// Ticker
	go ticker.InitState()

	// Orderbook
	go orderbook.InitState()

	// Trades
	go trades.InitState()

	// initi websocket client
	go wsclient.InitWebsocketClient("BTCUSD", "BTCUSD")
	go wsclient.InitWebsocketClient("ETHUSD", "ETHUSD")
	go wsclient.InitWebsocketClient("ETHBTC", "ETHBTC")

	// Set api routes
	router := mux.NewRouter()

	router.HandleFunc("/ticker", ticker.GetAll).Methods("GET")
	router.HandleFunc("/ticker/{pair}", ticker.Get).Methods("GET")

	router.HandleFunc("/orderbook", orderbook.GetAll).Methods("GET")
	router.HandleFunc("/orderbook/{pair}", orderbook.Get).Methods("GET")

	router.HandleFunc("/trades", trades.GetAll).Methods("GET")
	router.HandleFunc("/trades/{pair}", trades.Get).Methods("GET")

	// Start the server
	http.ListenAndServe(":8004", router)
}
