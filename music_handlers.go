package main

import (
	"encoding/json"
	"io"
	"log"
	"nb_proj3/apitypes"
	"nb_proj3/music"
	"net/http"
)

func addAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiAlbum

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := music.AddAlbum(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func removeAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiAlbumRemove

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := music.RemoveAlbum(deserialized.Token, deserialized.UUID)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func updateAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "PUT" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiAlbum

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := music.UpdateAlbum(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}
