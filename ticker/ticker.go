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

// Get a ticker from Kraken api
func getTicker(pair string) (aTickerResponse TickerResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%s", pair)

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
func write(pair string, result TickerResult) {
	ticker := Ticker{
		LastPrice: result.Last[0],
		BestBid:   result.Bid[0],
		BestAsk:   result.Ask[0],
		Datetime: time.Now().Unix(),
	}

	write := &writeOp{
		key:  pair,
		val:  ticker,
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
				tickerResponse, err := getTicker(xchK)
				if err == nil {
					write(sg3K, tickerResponse.Result[xchK])
				}
			}
		}(sg3Key, xchKey)
	}
}
