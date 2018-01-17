package ticker

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
	Ask       [3]string `json:"a"`
	Bid       [3]string `json:"b"`
	Last      [2]string `json:"c"`
	Volume    [2]string `json:"v"`
	Vwap      [2]string `json:"p"`
	NumTrades [2]int    `json:"t"`
	Low       [2]string `json:"l"`
	High      [2]string `json:"h"`
	Open      string    `json:"o"`
}

type TickerResponse struct {
	Error  []string                `json:"error"`
	Result map[string]TickerResult `json:"result"`
}

var pairs = map[string]string{
	"BTCUSD": "XXBTZUSD",
	"ETHUSD": "XETHZUSD",
	"LTCUSD": "XLTCZUSD"}

func getTicker(pair string) (aTickerResponse TickerResponse, err error) {
	httpCLI := &http.Client{
		Timeout: 10 * time.Second,
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

// GLOBAL VARS
var tickers = make(map[string]Ticker)

var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp)

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
func write(pair string) {
	ticker := Ticker{
		LastPrice: strconv.Itoa(rand.Intn(5)),
		BestBid:   strconv.Itoa(rand.Intn(5)),
		BestAsk:   strconv.Itoa(rand.Intn(5))}

	write := &writeOp{
		key:  pair,
		val:  ticker,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}

// Init service
func InitService() {

	for sg3Key, _ := range pairs {
		ticker := time.NewTicker(time.Millisecond * 1000)

		go func() {
			for range ticker.C {
				// tickerResponse, _ := getTicker(xchKey)
				//fmt.Printf("%s: %v\n", sg3Key, tickerResponse)
				write(sg3Key)
			}
		}()
	}
}

// HANDLERS
func GetTickerAll(w http.ResponseWriter, r *http.Request) {
	read := &readAllOp{resp: make(chan map[string]Ticker)}
	readsAll <- read
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(<-read.resp)
}

func GetTicker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	read := &readOneOp{
		key:  strings.ToUpper(params["pair"]),
		resp: make(chan Ticker)}
	readsOne <- read
	json.NewEncoder(w).Encode(<-read.resp)
}
