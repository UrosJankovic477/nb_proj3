package main

import (
	"encoding/json"
	"io"
	"log"
	"nb_proj3/apitypes"
	"nb_proj3/music"
	"nb_proj3/search"
	"net/http"
)

func getMultipleAlbumsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.Pagination

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	albums, status, err := music.GetMultypleAlbums(deserialized.Count, deserialized.Page)

	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	album_json, err := json.Marshal(albums)
	writer.Write([]byte(album_json))
	writer.WriteHeader(http.StatusOK)
}

func getAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiGetAlbum

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	album, status, err := music.GetAlbum(deserialized.Token, deserialized.AlbumUUID)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	album_json, err := json.Marshal(album)
	writer.Write([]byte(album_json))
	writer.WriteHeader(http.StatusOK)
}

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
	var deserialized apitypes.ApiAddAlbum

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	uuid, status, err := music.AddAlbum(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	json_uuid, _ := json.Marshal(uuid)
	writer.Write([]byte(json_uuid))
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

func getAlbumSongsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized string

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	songs, status, err := music.GetAlbumSongs(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	songs_json, _ := json.Marshal(songs)
	writer.Write([]byte(songs_json))
	writer.WriteHeader(http.StatusOK)
}

func getArtistsAlbumsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized string

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	albums, status, err := music.GetArtistsAlbums(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	album_json, _ := json.Marshal(albums)
	writer.Write([]byte(album_json))
	writer.WriteHeader(http.StatusOK)
}

func getPurchasedItemsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized string

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	items, status, err := music.GetPurchasedItems(deserialized)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	items_json, _ := json.Marshal(items)
	writer.Write([]byte(items_json))
	writer.WriteHeader(http.StatusOK)
}

func addSongToAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var deserialized apitypes.ApiAddSongToAlbum
	var song_bytes []byte
	mpr, err := reqptr.MultipartReader()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	json_part, err := mpr.NextPart()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(json_part)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	song_bytes_part, err := mpr.NextPart()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	song_bytes, err = io.ReadAll(song_bytes_part)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := music.AddSongToAlbum(deserialized.Token,
		deserialized.AlbumUUID,
		deserialized.SongData,
		song_bytes)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)

}

func removeSongFromAlbumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiRemoveSongFromAlbum

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := music.RemoveSongFromAlbum(deserialized.Token, deserialized.SongUUID)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)

}

func uploadAlbumArtHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var deserialized apitypes.ApiGetAlbum
	var art_bytes []byte
	mpr, err := reqptr.MultipartReader()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	json_part, err := mpr.NextPart()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(json_part)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	art_bytes_part, err := mpr.NextPart()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	art_bytes, err = io.ReadAll(art_bytes_part)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	status, err := music.UploadAlbumArt(deserialized.Token, deserialized.AlbumUUID, art_bytes)
	if err != nil {
		writer.WriteHeader(status)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func searchHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	var deserialized apitypes.ApiSearch

	err = json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	albums, status, err := search.SearchAlbums(deserialized.Query, deserialized.Genres, deserialized.MaxPrice, deserialized.Count, deserialized.Page)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	albums_json, err := json.Marshal(&albums)
	if err != nil {
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(albums_json)
}
