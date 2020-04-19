package main

import (
	"net/http"
	"io/ioutil"
)

func handleRootPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.URL.Path == "/" {
		body, err := ioutil.ReadFile("public/index.html")
		if err != nil {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(body))
		return
	}

	http.Error(w, "404 not found", http.StatusNotFound)
	return
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