package ticker

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	pairs := []string{"BTCUSD", "ETHUSD", "LTCUSD"}
	for _, pair := range pairs {
		ticker := time.NewTicker(time.Millisecond * 500)
		go func(p string) {
			for range ticker.C {
				write(p)
			}
		}(pair)
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
