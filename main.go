package main

import (
	"github.com/efrainmunoz/generic-api/ticker"
	"github.com/efrainmunoz/generic-api/orderbook"
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

	// Set api routes
	router := mux.NewRouter()
	router.HandleFunc("/ticker", ticker.GetAll).Methods("GET")
	router.HandleFunc("/ticker/{pair}", ticker.Get).Methods("GET")
	router.HandleFunc("/orderbook", orderbook.GetAll).Methods("GET")
	router.HandleFunc("/orderbook/{pair}", orderbook.Get).Methods("GET")

	// Start the server
	http.ListenAndServe(":8000", router)
}
