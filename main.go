package main

import (
	"github.com/efrainmunoz/generic-api/ticker"
	"github.com/gorilla/mux"
	"net/http"
)

// MAIN
func main() {

	// Ticker
	go ticker.InitState()
	go ticker.InitService()

	// Set api routes
	router := mux.NewRouter()
	router.HandleFunc("/ticker", ticker.GetTickerAll).Methods("GET")
	router.HandleFunc("/ticker/{pair}", ticker.GetTicker).Methods("GET")

	// Start the server
	http.ListenAndServe(":8000", router)
}
