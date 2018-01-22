package trades

// MODELS
type Trade struct {
	Price   string `json:"price"`
	Volume  string `json:"volume"`
	TradeAction string `json:"trade-action"`
}

type Trades struct {
	Trades    []Trade `json:"trades"`
	Timestamp int64   `json:"timestamp"`
}

type readAllOp struct {
	resp chan map[string]Trade
}

type readOneOp struct {
	key  string
	resp chan Trade
}

type writeOp struct {
	key  string
	val  Trade
	resp chan bool
}

// Models to deal with the third-party api responses
type IncomingTrade []interface{} // <price>, <volume>, <time>, <buy/sell>, <market/limit>, <miscellaneous>

type TradesResult map[string]interface{}
//type TradesResult struct {
//	Trades    []IncomingTrade
//	LastID    string `json:"last"`
//}

type TradesResponse struct {
	Error  []string     `json:"error"`
	Result TradesResult `json:"result"`
}
