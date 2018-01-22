package trades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GLOBAL VARS
var trades = make(map[string]Trades)
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

// Get trades from Kraken api
func getTrades(pair string) (aTradesResponse TradesResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	url := fmt.Sprintf("https://api.kraken.com/0/public/Trades?pair=%s", pair)

	// try to get kraken trades
	resp, err := httpCLI.Get(url)
	if err != nil {
		return TradesResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TradesResponse{}, err
	}

	// Unmarshal the json
	tradesResponse := TradesResponse{}
	err = json.Unmarshal(body, &tradesResponse)
	if err != nil {
		return TradesResponse{}, err
	}

	return tradesResponse, nil
}


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
			//if write.key == "BTCUSD" {
			//	fmt.Println(write.val)
			//}
			state[write.key] = write.val
			write.resp <- true
		}
	}
}

// WRITE new tickers
func write(sg3Pair string, xchPair string, result TradesResult) {
	if value, ok := result[xchPair].([]interface{}); ok {
		l := len(value)
		lastTrade :=  value[l-1]
		lastTradeSlc := lastTrade.([]interface{})
		var tradeAction string
		if lastTradeSlc[3].(string) == "s" {
			tradeAction = "sell"
		} else {
			tradeAction = "buy"
		}

		trade := Trade{
			Price: lastTradeSlc[0].(string),
			Volume: lastTradeSlc[1].(string),
			TradeAction: tradeAction}

		write := &writeOp{
			key:  sg3Pair,
			val: trade,
			resp: make(chan bool)}

		writes <- write
		<-write.resp
	}
}

// Init service
func InitService() {
	for sg3Key, xchKey := range pairs {
		go func(sg3K string, xchK string) {
			ticker := time.NewTicker(time.Millisecond * 1000)
			for range ticker.C {
				tradesResponse, err := getTrades(xchK)
				if err == nil {
					write(sg3K, xchK, tradesResponse.Result)
				}
			}
		}(sg3Key, xchKey)
	}
}
