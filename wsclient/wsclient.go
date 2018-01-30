package wsclient

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sort"
	"github.com/efrainmunoz/gem-api/ticker"
	"github.com/efrainmunoz/gem-api/trades"
	"github.com/efrainmunoz/gem-api/orderbook"
	"strconv"
)

// GLOBAL VARS

var pairs = map[string]string{
	"BTCUSD":  "BCTUSD",
	//"ETHUSD":  "ETHUSD",
	//"ETHBTC":  "ETHBTC",
}

type IncomingMessage struct {
	Type           string            `json:"type"`
	EventId        int               `json:"eventId"`
	SocketSequence int               `json:"socket_sequence"`
	Timestamp      int               `json:"timestamp"`
	Timestamps     int               `json:"timestamps"`
	Events         []json.RawMessage `json:"events"`
}

type TradeEvent struct {
	Type      string `json:"type"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	MakerSide string `json:"makerSide"`
}

type ChangeEvent struct {
	Type      string `json:"type"`
	Side      string `json:"side"`
	Price     string `json:"price"`
	Remaining string `json:"remaining"`
	Delta     string `json:"delta"`
	Reason    string `json:"reason"`
}

// sorting bids
type BidsByPrice []orderbook.Order
func (b BidsByPrice) Len()          int  { return len(b) }
func (b BidsByPrice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BidsByPrice) Less(i, j int) bool {
	l, _ := strconv.ParseFloat(b[i].Price, 64)
	r, _ := strconv.ParseFloat(b[j].Price, 64)
	return l > r
}

// sorting asks
type AsksByPrice []orderbook.Order
func (b AsksByPrice) Len()          int  { return len(b) }
func (b AsksByPrice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b AsksByPrice) Less(i, j int) bool {
	l, _ := strconv.ParseFloat(b[i].Price, 64)
	r, _ := strconv.ParseFloat(b[j].Price, 64)
	return l < r
}

// Init service
func InitWebsocketClient(sg3, xchg string) {
				// Get initial ticker
				tickerInitialResponse, _ := ticker.GetTicker(xchg)
				ticker.Write(sg3, tickerInitialResponse)

				// Last price
				lastPrice := tickerInitialResponse.Last

				// slices for the orderbook
				bids := make([]orderbook.Order, 0)
				asks := make([]orderbook.Order, 0)

				// make a websocket connection
				fmt.Printf("Attempting conn to for %s\n", xchg)
				url := fmt.Sprintf("wss://api.gemini.com/v1/marketdata/%s", xchg)
				conn, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					fmt.Println(err)
				}

				// close connection in case of unexpected function exit
				defer conn.Close()

				// listen for incoming gemini messages
				for {
					// try to read a message
					fmt.Printf("Reading msg from %s conn\n", xchg)
					_, msg, err := conn.ReadMessage()
					if err != nil {
						fmt.Println(err)
						return
					}

					// Parse incoming JSON
					incomingMessage := IncomingMessage{}
					err = json.Unmarshal(msg, &incomingMessage)
					if err != nil {
						panic(err)
					}

					// manage initial response
					if incomingMessage.SocketSequence == 0 {
						// iterate over the events array in the incomingMessage
						for _, event := range incomingMessage.Events {
							// parse event
							initial := ChangeEvent{}
							err := json.Unmarshal(event, &initial)
							if err != nil {
								panic(err)
							}

							// append bid/ask to in-memory slice
							if initial.Side == "bid" {
								bids = append([]orderbook.Order{
									orderbook.Order{Price: initial.Price, Volume: initial.Remaining}},
									bids...
								)
							} else {
								asks = append(asks, orderbook.Order{Price: initial.Price, Volume: initial.Remaining})
							}
						}
					}

					// if not initial response
					if incomingMessage.SocketSequence > 0 {
						// iterate over the events array in the incomingMessage
						for _, event := range incomingMessage.Events {
							// parse event into a map of interfaces to be able to
							// access the values via indexing
							var rawEvent map[string]interface{}
							err := json.Unmarshal(event, &rawEvent)
							if err != nil {
								fmt.Println(err)
							} else {
								// if event is a trade do something
								if rawEvent["type"] == "trade" {
									trade := TradeEvent{}
									err := json.Unmarshal(event, &trade)
									if err != nil {
										fmt.Println(err)
									} else {
										lastPrice = trade.Price
										// Propagate new ticker
										ticker.Write(sg3, ticker.TickerResponse{
											Bid: bids[0].Price, Ask: asks[0].Price, Last: lastPrice})
										// Propagate new trade
										var tradeAction string
										if trade.MakerSide == "ask" {
											tradeAction = "buy"
										} else {
											tradeAction = "sell"
										}
										trades.Write(sg3, trades.Trade{
											Price:       trade.Price,
											Volume:      trade.Amount,
											TradeAction: tradeAction,
										})
									}
									// if event is a chage to the orderbook do something
								} else if rawEvent["type"] == "change" {
									change := ChangeEvent{}
									err := json.Unmarshal(event, &change)
									if err != nil {
										fmt.Println(err)
									} else {
										// if price is 0 remove entry from slice
										if change.Remaining == "0" {
											if change.Side == "bid" {
												for i, v := range bids {
													if v.Price == change.Price {
														bids = append(bids[:i], bids[i+1:]...)
														break
													}
												}
											} else {
												for i, v := range asks {
													if v.Price == change.Price {
														asks = append(asks[:i], asks[i+1:]...)
														break
													}
												}
											}
											// if price is not 0
										} else {
											if change.Side == "bid" {
												found := false
												for i, v := range bids {
													// if price point found in the slice update the Oder
													if v.Price == change.Price {
														bids[i] = orderbook.Order{
															Price: change.Price,
															Volume: change.Remaining,
														}
														//fmt.Printf("Bids update: %s\n\n\n", bids)
														found = true
														break
													}
												}
												// if price point not found append the value to the slice and sort it
												if found == false {
													bids = append(bids, orderbook.Order{
														Price: change.Price,
														Volume: change.Remaining,
													})
													sort.Sort(BidsByPrice(bids))
												}
												//fmt.Println(bids)
											} else {
												found := false
												for i, v := range asks {
													// if price point found in the slice update the Order
													if v.Price == change.Price {
														asks[i] = orderbook.Order{
															Price: change.Price,
															Volume: change.Remaining,
														}
														//fmt.Printf("Asks update: %s\n\n\n", asks)
														found = true
														break
													}
												}
												// if price point not found append the value to the slice and sort it
												if found == false {
													asks = append(asks,
														orderbook.Order{Price: change.Price, Volume: change.Remaining})
													sort.Sort(AsksByPrice(asks))
													//fmt.Printf("Asks add: %s\n\n\n", asks)
												}
											}
										}
									}

									// Propagate new ticker
									newTicker := ticker.TickerResponse{
										Bid:  bids[0].Price,
										Ask:  asks[0].Price,
										Last: lastPrice,
									}
									ticker.Write(sg3, newTicker)

									// Prpagate new Orderbook
									newOrderbook := orderbook.OrderbookResponse{bids, asks}
									orderbook.Write(sg3, newOrderbook)
								}
							}
						}
					}
				}
}