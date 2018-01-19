package trades

// MODELS
type Trade struct {
	Price   string `json:"price"`
	Volume  string `json:"volume"`
	BuySell string `json:"buy-sell"`
}

type Trades struct {
	Trades    []Trade `json:"trades"`
	Timestamp int64   `json:"timestamp"`
}

type readAllOp struct {
	resp chan map[string]Trades
}

type readOneOp struct {
	key  string
	resp chan Trades
}

type writeOp struct {
	key  string
	val  Trades
	resp chan bool
}

// Models to deal with the third-party api responses
type IncomingTrade []interface{} // <price>, <volume>, <time>, <buy/sell>, <market/limit>, <miscellaneous>

type TradesResult struct {
	Trades    []IncomingTrade `json:"asks"`
	Bids      []IncomingTrade `json:"bids"`
}

type TradesResponse struct {
	Error  []string                   `json:"error"`
	Result map[string]TradesResult `json:"result"`
}
