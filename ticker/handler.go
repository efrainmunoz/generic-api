package ticker

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strings"
)

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