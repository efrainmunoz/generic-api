package orderbook

// MODELS
type Order struct {
	Price    string
	Volume   string
}

type Orderbook struct {
	Asks []Order
	Bids []Order
	Timestamp int64
}

type readAllOp struct {
	resp chan map[string]Orderbook
}

type readOneOp struct {
	key  string
	resp chan Orderbook
}

type writeOp struct {
	key  string
	val  Orderbook
	resp chan bool
}

// Models to deal with the third-party api responses
type Ask []interface{} // <price>, <volume>, <timestamp>
type Bid []interface{} // <price>, <volume>, <timestamp>

type OrderbookResult struct {
	Asks      []Ask `json:"asks"`
	Bids      []Bid `json:"bids"`
}

type OrderbookResponse struct {
	Error  []string                   `json:"error"`
	Result map[string]OrderbookResult `json:"result"`
}