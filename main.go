package main

import (
	"github.com/efrainmunoz/generic-api/orderbook"
	"github.com/efrainmunoz/generic-api/ticker"
	"github.com/efrainmunoz/generic-api/trades"
	"github.com/gorilla/mux"
	"net/http"
)

// MAIN
func main() {

	// Ticker
	go ticker.InitState()
	go ticker.InitService()

	// Orderbook
	go orderbook.InitState()
	go orderbook.InitService()

	// Trades
	go trades.InitState()
	go trades.InitService()

	// Set api routes
	router := mux.NewRouter()
	router.HandleFunc("/ticker", ticker.GetAll).Methods("GET")
	router.HandleFunc("/ticker/{pair}", ticker.Get).Methods("GET")
	router.HandleFunc("/orderbook", orderbook.GetAll).Methods("GET")
	router.HandleFunc("/orderbook/{pair}", orderbook.Get).Methods("GET")
	router.HandleFunc("/trades", trades.GetAll).Methods("GET")
	router.HandleFunc("/trades/{pair}", trades.Get).Methods("GET")

	// Start the server
	http.ListenAndServe(":8001", router)
}
