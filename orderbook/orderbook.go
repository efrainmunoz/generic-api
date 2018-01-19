package orderbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GLOBAL VARS
var orderbook = make(map[string]Orderbook)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)
var pairs = map[string]string{
	"BTCUSD": "XXBTZUSD",
	"ETHUSD": "XETHZUSD",
	"ETHBTC": "XETHXXBT",
	"LTCUSD": "XLTCZUSD",
	"LTCBTC": "XLTCXXBT",
	"XRPUSD": "XXRPZUSD",
	"XRPBTC": "XXRPXXBT",
	"ZECUSD": "XZECZUSD",
	"ZECBTC": "XZECXXBT",
	"XMRUSD": "XXMRZUSD",
	"XMRBTC": "XXMRXXBT",
	"DASHUSD": "DASHUSD",
	"DASHBTC": "DASHXBT",
	"BCHUSD": "BCHUSD",
	"BCHBTC": "BCHXBT",
	"ETCUSD": "XETCZUSD",
	"ETCBTC": "XETCXXBT",
}

// Get an orderbook from Kraken api
func getOrderbook(pair string) (aOrderbookResponse OrderbookResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("https://api.kraken.com/0/public/Depth?pair=%s", pair)

	// try to get kraken orderbook
	resp, err := httpCLI.Get(url)
	if err != nil {
		return OrderbookResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OrderbookResponse{}, err
	}

	// Unmarshal the json
	orderbookResponse := OrderbookResponse{}
	err = json.Unmarshal(body, &orderbookResponse)
	if err != nil {
		return OrderbookResponse{}, err
	}

	return orderbookResponse, nil
}


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
func write(pair string, result OrderbookResult) {
	var asks []Order
	var bids []Order

	for _, ask := range result.Asks {
		order := Order{
			Price: ask[0].(string),
			Volume: ask[1].(string)}
		asks = append(asks, order)
	}

	for _, bid := range result.Bids {
		bid :=  Order{
			Price: bid[0].(string),
			Volume: bid[1].(string)}
		bids = append(asks, bid)
	}

	orderbook := Orderbook{
		Asks: asks,
		Bids: bids,
		Timestamp: time.Now().Unix()}

	write := &writeOp{
		key:  pair,
		val:  orderbook,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}

// Init service
func InitService() {
	for sg3Key, xchKey := range pairs {
		go func(sg3K string, xchK string) {
			ticker := time.NewTicker(time.Millisecond * 1000)
			for range ticker.C {
				orderbookResponse, err := getOrderbook(xchK)
				if err == nil {
					write(sg3K, orderbookResponse.Result[xchK])
				}
			}
		}(sg3Key, xchKey)
	}
}
