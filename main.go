package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Client struct {
	Id int
	Phone string
	Pin string
}

type NewClientPayload struct {
	Phone string
	Pin string
	Retyped_pin string
}

type ErrorPayload struct {
	Error bool `json:"error"`
	ErrorCode string `json:"error_code"`
}

var HumbleClientsStorage []Client

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
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	var payload NewClientPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	for _, v := range HumbleClientsStorage {
		if v.Phone == payload.Phone {
			js, err := json.Marshal(ErrorPayload{ true, "NUMBER_ALREADY_REGISTERED" })
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(js)
			return
		}
	}

	newClientId := len(HumbleClientsStorage) + 1
	client := Client{ newClientId, payload.Phone, payload.Pin }
	HumbleClientsStorage = append(HumbleClientsStorage, client)

	js, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	return
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