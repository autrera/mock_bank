package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	b64 "encoding/base64"
	"strconv"
)

type Client struct {
type Transfer struct {
	Id int
	Amount int
	ClientId int
	CreatedBy int
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
	Error bool `json:"error"`
	ClientId int `json:"client_id"`
	Token string `json:"token"`
}

type JsonResponse struct {
	Payload interface{}
}

type ErrorPayload struct {
	Error bool `json:"error"`
	ErrorCode string `json:"error_code"`
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