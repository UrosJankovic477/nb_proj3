package music

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/apiupload"
	"nb_proj3/connection"
	"nb_proj3/users"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMultypleAlbums(count int64, page int64) ([]apitypes.Album, int, error) {
	db := connection.GetDatabase()
	skip := page * count
	opts := options.FindOptions{Limit: &count, Skip: &skip}
	cursor, err := db.Collection("items").Find(context.TODO(), bson.D{{}}, &opts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	results := make([]apitypes.Album, 0)
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return results, http.StatusOK, nil
}

func GetAlbum(token string, album_uuid string) (apitypes.AlbumData, int, error) {
	db := connection.GetDatabase()
	owns := false
	username := ""

	if token != "" {
		user, _, err := users.GetUser(token)
		if err == nil {
			username = user.Username
		}
		err = db.Collection("users").FindOne(context.TODO(), bson.D{{"_id", user.Username}, {"purchaseditems", album_uuid}}).Err()
		if err == nil {
			owns = true
		}
	}

	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"_id", album_uuid}})
	var album apitypes.Album
	err := result.Decode(&album)
	if err != nil {
		return apitypes.AlbumData{}, http.StatusNotFound, err
	}
	cursor, err := db.Collection("songs").Find(context.TODO(), bson.D{{"albumuuid", album_uuid}})
	songs := make([]apitypes.Song, 0)
	err = cursor.All(context.TODO(), &songs)
	if err != nil {
		return apitypes.AlbumData{}, http.StatusNotFound, err
	}
	created := username == album.ArtistUsername

	album_data := apitypes.AlbumData{
		UUID:           album.UUID,
		Title:          album.Title,
		ArtistUsername: album.ArtistUsername,
		Songs:          songs,
		Price:          album.Price,
		Created:        created,
		Owns:           owns,
	}
	return album_data, http.StatusOK, err
}

func AddAlbum(api_album apitypes.ApiAddAlbum) (string, int, error) {
	if api_album.Title == "" {
		return "", http.StatusBadRequest, errors.New("Album title cannot be empty")
	}
	user, status, err := users.GetUser(api_album.Token)
	if err != nil {
		return "", status, err
	}
	var album_uuid string
	album_uuid = uuid.New().String()
	album := apitypes.Album{Title: api_album.Title, UUID: album_uuid}
	album.ArtistUsername = user.Username
	album.Price = 0
	album.UUID = album_uuid
	db := connection.GetDatabase()
	db.Collection("albums").InsertOne(context.TODO(), album)
	return album_uuid, http.StatusOK, nil
}

func RemoveAlbum(token string, uuid string) (int, error) {
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	username := user.Username
	db := connection.GetDatabase()
	_, err = db.Collection("albums").DeleteOne(context.TODO(), bson.D{{"_id", uuid}, {"artistusername", username}})

	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = db.Collection("albums").DeleteMany(context.TODO(), bson.D{{"albumuuid", uuid}})

	return http.StatusOK, nil
}

func UpdateAlbum(api_album apitypes.ApiAlbum) (int, error) {
	user, status, err := users.GetUser(api_album.Token)
	if err != nil {
		return status, err
	}
	username := user.Username
	db := connection.GetDatabase()
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"_id", api_album.UUID}, {"artistusername", username}})
	var old_album apitypes.Album
	result.Decode(&old_album)
	new_album := api_album.Album
	if new_album.Title == "" {
		new_album.Title = old_album.Title
	}

	_, err = db.Collection("albums").UpdateOne(context.TODO(), bson.D{{"_id", old_album.UUID}, {"artistusername", username}}, new_album)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func GetAlbumsFromArtist(username string, limit int64, page int64) ([]apitypes.Album, int, error) {
	db := connection.GetDatabase()
	opts := options.Find().SetLimit(limit).SetSkip((page - 1) * limit)
	cursor, err := db.Collection("albums").Find(context.TODO(), bson.D{{"artistusername", username}}, opts)
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

func AddSongToAlbum(token string, album_uuid string, song apitypes.SongData, song_bytes []byte) (int, error) {
	if song.Title == "" {
		return http.StatusBadRequest, errors.New("Title cannot be empty")
	}
	if song.Genre == "" {
		return http.StatusBadRequest, errors.New("Genre cannot be empty")
	}
	if song.Price < 0 {
		return http.StatusBadRequest, errors.New("Price cannot be negative")
	}

	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	result := db.Collection("albums").FindOne(context.TODO(),
		bson.D{{"_id", album_uuid}, {"artistusername", user.Username}})
	var album apitypes.Album
	err = result.Decode(&album)
	if err != nil {
		return http.StatusBadRequest, err
	}
	song_uuid := uuid.New().String()
	length, err := apiupload.UploadAudio(song_uuid, song_bytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	added_song := apitypes.Song{
		AlbumUUID: album_uuid,
		UUID:      song_uuid,
		Title:     song.Title,
		Genre:     song.Genre,
		Length:    length,
		Price:     song.Price,
	}

	db.Collection("songs").InsertOne(context.TODO(), added_song)

	album.Songs = append(album.Songs, apitypes.SongData{
		UUID:  song_uuid,
		Title: song.Title,
		Price: song.Price,
		Genre: song.Genre,
	})
	album.Price += song.Price
	_, err = db.Collection("albums").ReplaceOne(context.TODO(), bson.D{{"_id", album_uuid}}, album)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func RemoveSongFromAlbum(token string, song_uuid string) (int, error) {
	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}

	result := db.Collection("songs").FindOneAndDelete(context.TODO(), bson.D{{"_id", song_uuid}})
	var song apitypes.Song
	err = result.Decode(&song)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	result = db.Collection("albums").FindOne(context.TODO(),
		bson.D{{"_id", song.AlbumUUID}, {"artistusername", user.Username}})

	var album apitypes.Album
	result.Decode(&album)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, err
	}

	songs := make([]apitypes.SongData, 0)
	for _, v := range album.Songs {
		if v.UUID == song_uuid {
			continue
		}
		songs = append(songs, v)
	}
	apiupload.DeleteAudio(song_uuid)
	apiupload.DeleteAudio("pre_" + song_uuid)
	db.Collection("songs").DeleteOne(context.TODO(), bson.D{{"_id", song_uuid}})
	album.Songs = songs
	album.Price -= song.Price
	db.Collection("albums").ReplaceOne(context.TODO(), bson.D{{"_id", album.UUID}}, album)
	return http.StatusOK, nil
}

func UploadAlbumArt(token string, album_uuid string, art_bytes []byte) (int, error) {
	db := connection.GetDatabase()
	user, status, err := users.GetUser(token)
	if err != nil {
		return status, err
	}
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"_id", album_uuid}, {"ArtistUsername", user.Username}})
	var album apitypes.Album
	result.Decode(&album)

	apiupload.DeleteImage(album_uuid)
	apiupload.UploadImage(album_uuid, art_bytes)

	return http.StatusOK, nil
}
