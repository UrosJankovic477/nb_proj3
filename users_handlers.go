package main

import (
	"encoding/json"
	"io"
	"log"
	"nb_proj3/apitypes"
	"nb_proj3/users"
	"net/http"
)

func registerHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	user := apitypes.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = users.CreateUser(user.Username, user.PasswordHash)
	if err != nil {
		writer.WriteHeader(http.StatusConflict)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func loginHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	user := apitypes.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken, status, err := users.Login(user.Username, user.PasswordHash)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}

	token_json, err := json.Marshal(&sessionToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	writer.Write(token_json)
	writer.WriteHeader(http.StatusOK)
}

func addToCartHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	item := apitypes.ApiAddToCart{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := users.AddToCart(item.Token, item.Buyable)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func removeItemFromCartHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	item := apitypes.ApiRemoveFromCart{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := users.RemoveFromCart(item.Token, item.ItemUUID)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func checkoutHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	token := ""
	err = json.Unmarshal(body, &token)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := users.Checkout(token)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}
