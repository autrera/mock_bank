package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	// "fmt"
)

type Client struct {
	Id int `json:"id,string"`
	Phone string `json:"phone"`
	Pin string `json:"pin"`
}

type Transfer struct {
	Id        int `json:"id"`
	Amount    int `json:"amount"`
	ClientId  int `json:"client_id"`
	CreatedBy int `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type NewClientRequestPayload struct {
	Phone string
	Pin string
	Retyped_pin string
}

type NewClientResponsePayload struct {
	Error bool `json:"error"`
	ClientId int `json:"client_id"`
	Token string `json:"token"`
}

type NewLoginRequestPayload struct {
	Phone string
	Pin string
}

type NewLoginResponsePayload struct {
	Error    bool   `json:"error"`
	ClientId int    `json:"client_id"`
	Token    string `json:"token"`
}

type NewTransferRequestPayload struct {
	Amount int `json:"amount,string"`
	Phone  string `json:"phone"`
	Pin    string `json:"pin"`
}

type JsonResponse struct {
	Payload interface{}
}

type ErrorPayload struct {
	Error bool `json:"error"`
	ErrorCode string `json:"error_code"`
}

type BalancePayload struct {
	Error bool `json:"error"`
	Balance int `json:"balance"`
	Transfers []Transfer `json:"transfers"` 
}

var HumbleClientsStorage []Client
var HumbleTransfersStorage []Transfer

func sendJsonResponse(w http.ResponseWriter, response JsonResponse, httpCode int) {
	js, err := json.Marshal(response.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(js)
	return
}

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
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	var payload NewLoginRequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
	}

	for _, v := range HumbleClientsStorage {
		if v.Phone == payload.Phone && v.Pin == payload.Pin {
			mockJson := `{ "id":"` + strconv.Itoa(v.Id) + `", "phone":"` + v.Phone + `", "pin":"` + v.Pin + `" }`
			token := b64.URLEncoding.EncodeToString([]byte(mockJson))
			sendJsonResponse(w, JsonResponse{ NewLoginResponsePayload{ false, v.Id, token } }, 200)
			return
		}
	}

	sendJsonResponse(w, JsonResponse{ ErrorPayload{ true, "UNAUTHORIZED" }}, 403)
	return
}

func handleNewClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	var payload NewClientRequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	for _, v := range HumbleClientsStorage {
		if v.Phone == payload.Phone {
			sendJsonResponse(w, JsonResponse{ ErrorPayload{ true, "NUMBER_ALREADY_REGISTERED" }}, 400)
			return
		}
	}

	newClientId := len(HumbleClientsStorage) + 1
	client := Client{ newClientId, payload.Phone, payload.Pin }
	HumbleClientsStorage = append(HumbleClientsStorage, client)

	newClientTransferId := len(HumbleTransfersStorage) + 1
	transfer := Transfer{ newClientTransferId, 1000000, newClientId, 0 }
	HumbleTransfersStorage = append(HumbleTransfersStorage, transfer)

	mockJson := `{ "id":"` + strconv.Itoa(client.Id) + `", "phone":"` + client.Phone + `", "pin":"` + client.Pin + `" }`
	token := b64.URLEncoding.EncodeToString([]byte(mockJson))
	sendJsonResponse(w, JsonResponse{ NewClientResponsePayload{ false, client.Id, token } }, 200)
	return
}

func handleNewTransfers(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	decodedString, _ := b64.URLEncoding.DecodeString(token)

	var client Client
	err := json.NewDecoder(strings.NewReader(string(decodedString))).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newTransferRequest NewTransferRequestPayload
	err = json.NewDecoder(r.Body).Decode(&newTransferRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clientToReceiveTransfer Client
	for _, v := range HumbleClientsStorage {
		if v.Phone == newTransferRequest.Phone {
			clientToReceiveTransfer = v
		}
	}

	if clientToReceiveTransfer.Phone == "" {
	}

	newTransferId := len(HumbleTransfersStorage) + 1
	newTransfer := Transfer{newTransferId, newTransferRequest.Amount * 100, clientToReceiveTransfer.Id, client.Id, time.Now().Format(time.UnixDate)}
	HumbleTransfersStorage = append(HumbleTransfersStorage, newTransfer)

	sendJsonResponse(w, JsonResponse{newTransfer}, 200)
	return
}

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	decodedString, _ := b64.URLEncoding.DecodeString(token)

	var client Client
	err := json.NewDecoder(strings.NewReader(string(decodedString))).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	amount := 0
	transfers := []Transfer{}
	for _, v := range HumbleTransfersStorage {
		if v.ClientId == client.Id || v.CreatedBy == client.Id {
			if v.CreatedBy == client.Id {
				v.Amount = v.Amount * -1
			}
			v.Amount = v.Amount / 100
			amount += v.Amount
			transfers = append(transfers, v)
		}
	}

	sendJsonResponse(w, JsonResponse{BalancePayload{false, amount, transfers}}, 200)
	return
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
