package ticker

// MODELS
type Ticker struct {
	LastPrice string `json:"lastprice"`
	BestBid   string `json:"bestbid"`
	BestAsk   string `json:"bestask"`
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

type TickerResult struct {
	Ask       [3]string `json:"a"` // <price>, <whole lot volume>, <lot volume>
	Bid       [3]string `json:"b"` // <price>, <whole lot volume>, <lot volume>
	Last      [2]string `json:"c"` // <price>, <lot volume>
	Volume    [2]string `json:"v"` // <today>, <last 24 hours>
	Vwap      [2]string `json:"p"` // <today>, <last 24 hours>
	NumTrades [2]int    `json:"t"` // <today>, <last 24 hours>
	Low       [2]string `json:"l"` // <today>, <last 24 hours>
	High      [2]string `json:"h"` // <today>, <last 24 hours>
	Open      string    `json:"o"` // today's opening price
}

type TickerResponse struct {
	Error  []string                `json:"error"`
	Result map[string]TickerResult `json:"result"`
}