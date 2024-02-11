package main

import (
	"context"
	"encoding/json"
	"io"
	"nb_proj3/apitypes"
	"nb_proj3/apiupload"
	"nb_proj3/connection"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	connection.Init()

	router := http.NewServeMux()

	router.HandleFunc("/api/register", registerHandler)
	router.HandleFunc("/api/login", loginHandler)
	router.HandleFunc("/api/addSongToCart", addSongToCartHandler)
	router.HandleFunc("/api/addAlbumToCart", addAlbumToCartHandler)
	router.HandleFunc("/api/removeItemFromCart", removeItemFromCartHandler)
	router.HandleFunc("/api/checkout", checkoutHandler)
	router.HandleFunc("/api/emptyCart", emptyCartHandler)
	router.HandleFunc("/api/album", getAlbumHandler)
	router.HandleFunc("/api/addAlbum", addAlbumHandler)
	router.HandleFunc("/api/removeAlbum", removeAlbumHandler)
	router.HandleFunc("/api/updateAlbum", updateAlbumHandler)
	router.HandleFunc("/api/getMultipleAlbums", getMultipleAlbumsHandler)
	router.HandleFunc("/api/getLoggedInUsersData", getLoggedInUsersDataHandler)
	router.HandleFunc("/api/getAlbumSongs", getAlbumHandler)
	router.HandleFunc("/api/getArtistsAlbums", getArtistsAlbumsHandler)
	router.HandleFunc("/api/getPurchasedItems", getPurchasedItemsHandler)
	router.HandleFunc("/api/addSongToAlbum", addSongToAlbumHandler)
	router.HandleFunc("/api/removeSongFromAlbum", removeSongFromAlbumHandler)
	router.HandleFunc("/api/uploadAlbumArt", uploadAlbumArtHandler)
	router.HandleFunc("/api/getCart", getCartHandler)

	router.HandleFunc("/album/", serveAlbumHandler)
	router.HandleFunc("/artist/", serveArtistHandler)
	router.HandleFunc("/previews/", previewSongHandler)
	router.HandleFunc("/download/", downloadSongHandler)
	router.HandleFunc("/albumart/", albumArtHandler)

	router.HandleFunc("/api/search", searchHandler)
	router.HandleFunc("/api/getSearchOptions", getSearchOptions)

	router.Handle("/", http.FileServer(http.Dir("wwwroot")))

	http.ListenAndServe("localhost:8080", router)
}

func serveAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	http.ServeFile(writer, reqptr, "./wwwroot/album.html")
}

func serveArtistHandler(writer http.ResponseWriter, reqptr *http.Request) {
	http.ServeFile(writer, reqptr, "./wwwroot/artist.html")
}

func previewSongHandler(writer http.ResponseWriter, reqptr *http.Request) {
	path := strings.TrimPrefix(reqptr.URL.Path, "/previews/")
	path = "pre_" + path
	reader, err := apiupload.DownloadAudio(path)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(writer, reader)
}

func albumArtHandler(writer http.ResponseWriter, reqptr *http.Request) {
	path := strings.TrimPrefix(reqptr.URL.Path, "/albumart/")
	reader, err := apiupload.DownloadImage(path)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	img, err := io.ReadAll(reader)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	content_type := http.DetectContentType(img)
	writer.Header().Set("Content-Type", content_type)
	writer.Write(img)
}

func getSearchOptions(writer http.ResponseWriter, reqptr *http.Request) {
	db := connection.GetDatabase()
	genres, err := db.Collection("songs").Distinct(context.TODO(), "genre", bson.D{})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	opts := options.FindOne().SetSort(bson.D{{"price", -1}}).SetProjection(bson.D{{"price", 1}})
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{}, opts)
	var decoded struct {
		Price float64
	}
	var max_price float64
	err = result.Decode(&decoded)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	max_price = decoded.Price
	var search_options apitypes.ApiSearchOptions
	genres_string := make([]string, 0)
	for _, v := range genres {
		genres_string = append(genres_string, v.(string))
	}
	search_options.Genres = genres_string
	search_options.MaxPrice = max_price

	search_options_json, err := json.Marshal(&search_options)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(search_options_json))
}

func downloadSongHandler(writer http.ResponseWriter, reqptr *http.Request) {
	path := strings.TrimPrefix(reqptr.URL.Path, "/download/")
	reader, err := apiupload.DownloadAudio(path)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(writer, reader)
}
