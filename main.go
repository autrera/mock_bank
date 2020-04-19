package main

import (
	"net/http"
)

func handleRootPath(w http.ResponseWriter, r *http.Request) {
}

func handleNewLogin(w http.ResponseWriter, r *http.Request) {
}

func handleNewClient(w http.ResponseWriter, r *http.Request) {
}

func handleNewTransfers(w http.ResponseWriter, r *http.Request) {
}

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/", handleRootPath)
	http.HandleFunc("/login", handleNewLogin)
	http.HandleFunc("/clients", handleNewClient)
	http.HandleFunc("/transfers", handleNewTransfers)
	http.HandleFunc("/balance", handleGetBalance)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}