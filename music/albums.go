package music

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"nb_proj3/users"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAlbum(api_album apitypes.ApiAlbum) (int, error) {
	user, status, err := users.GetUser(api_album.Title)
	if err != nil {
		return status, err
	}
	album := api_album.Album
	album.ArtistUsername = user.Username
	album.UUID = uuid.New().String()
	db := connection.GetDatabase()
	db.Collection("albums").InsertOne(context.TODO(), album)
	return http.StatusOK, nil
}

func RemoveAlbum(token string, uuid string) (int, error) {
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	username := user.Username
	db := connection.GetDatabase()
	_, err = db.Collection("albums").DeleteOne(context.TODO(), bson.D{{"UUID", uuid}, {"ArtistUsername", username}})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func UpdateAlbum(api_album apitypes.ApiAlbum) (int, error) {
	user, status, err := users.GetUser(api_album.Token)
	if err != nil {
		return status, err
	}
	username := user.Username
	db := connection.GetDatabase()
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"UUID", api_album.UUID}, {"ArtistUsername", username}})
	var old_album apitypes.Album
	result.Decode(&old_album)
	new_album := api_album.Album
	if new_album.Title == "" {
		new_album.Title = old_album.Title
	}
	if new_album.AlbumArtLink == "" {
		new_album.AlbumArtLink = old_album.AlbumArtLink
	}
	new_album.Price = old_album.Price
	new_album.Songs = old_album.Songs

	_, err = db.Collection("albums").UpdateOne(context.TODO(), bson.D{{"UUID", old_album.UUID}, {"ArtistUsername", username}}, new_album)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func GetAlbumsFromArtist(username string, limit int64, page int64) ([]apitypes.Album, int, error) {
	db := connection.GetDatabase()
	opts := options.Find().SetLimit(limit).SetSkip((page - 1) * limit)
	cursor, err := db.Collection("albums").Find(context.TODO(), bson.D{{"ArtistUsername", username}}, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var albums []apitypes.Album
	err = cursor.All(context.TODO(), &albums)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return albums, http.StatusOK, nil
}

func AddSongToAlbum(token string, album_uuid string, song apitypes.Song) (int, error) {
	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	result := db.Collection("albums").FindOne(context.TODO(),
		bson.D{{"UUID", album_uuid}, {"ArtistUsername", user.Username}})

	var album apitypes.Album
	result.Decode(&album)
	album.Songs = append(album.Songs, song)
	_, err = db.Collection("albums").UpdateByID(context.TODO(), album_uuid, album)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	album.Price += song.Price
	return http.StatusOK, nil
}

func RemoveSongFromAlbum(token string, album_uuid string, song_uuid string) (int, error) {
	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	result := db.Collection("albums").FindOne(context.TODO(),
		bson.D{{"UUID", album_uuid}, {"ArtistUsername", user.Username}})

	var album apitypes.Album
	result.Decode(&album)

	idx := -1
	var removed apitypes.Song
	for i, v := range album.Songs {
		if v.UUID == song_uuid {
			removed = v
			idx = i
			break
		}
	}
	if idx == -1 {
		return http.StatusNotFound, errors.New("Song not found")
	}
	album.Songs = append(album.Songs[:idx], album.Songs[idx+1:]...)

	_, err = db.Collection("albums").UpdateByID(context.TODO(), album_uuid, album)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	album.Price -= removed.Price
	return http.StatusOK, nil
}

func SetAlbumArt(token string, album_uuid string, file_path string) (int, error) {
	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"UUID", album_uuid}, {"ArtistUsername", user.Username}})
	var album apitypes.Album
	result.Decode(&album)
	album.AlbumArtLink = file_path
	return http.StatusOK, nil
}
