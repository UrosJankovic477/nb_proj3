package main

import (
	"nb_proj3/connection"
	"net/http"
)

func main() {

	connection.Init()

	router := http.NewServeMux()

	router.HandleFunc("/api/register", registerHandler)
	router.HandleFunc("/api/login", loginHandler)
	router.HandleFunc("/api/addAlbum", addAlbumHandler)
	router.HandleFunc("/api/removeAlbum", removeAlbumHandler)
	router.HandleFunc("/api/updateAlbum", updateAlbumHandler)
	//TODO: Search engine

	router.Handle("/", http.FileServer(http.Dir("wwwroot")))

	http.ListenAndServe("localhost:8080", router)
}
