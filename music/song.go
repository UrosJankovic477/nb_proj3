package music

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"nb_proj3/users"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAlbumSongs(album_uuid string) ([]apitypes.Song, int, error) {
	db := connection.GetDatabase()
	cursor, err := db.Collection("songs").Find(context.TODO(), bson.D{{"albumuuid", album_uuid}})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	songs := make([]apitypes.Song, 0)
	cursor.All(context.TODO(), &songs)

	return songs, http.StatusOK, nil
}

func GetArtistsAlbums(artist_username string) ([]apitypes.ApiSearchResult, int, error) {
	db := connection.GetDatabase()
	cursor, err := db.Collection("albums").Find(context.TODO(), bson.D{{"artistusername", artist_username}})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	albums := make([]apitypes.Album, 0)
	cursor.All(context.TODO(), &albums)

	items := make([]apitypes.ApiSearchResult, 0)
	for _, v := range albums {
		user, _, _ := users.GetUserFromUsername(v.ArtistUsername)

		items = append(items, apitypes.ApiSearchResult{
			UUID:              v.UUID,
			Title:             v.Title,
			ArtistUsername:    v.ArtistUsername,
			ArtistDisplayname: user.DisplayName,
			Price:             v.Price,
		})
	}

	return items, http.StatusOK, nil
}

func GetPurchasedItems(token string) ([]apitypes.Album, int, error) {
	user, status, err := users.GetUser(token)
	if err != nil {
		return nil, status, err
	}
	db := connection.GetDatabase()
	opts := options.FindOne().SetProjection(bson.D{{"purchaseditems", 1}})
	result := db.Collection("users").FindOne(context.TODO(), bson.D{{"_id", user.Username}}, opts)
	var temp apitypes.TempAlbums
	err = result.Decode(&temp)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return temp.Albums, http.StatusOK, nil
}

func UpdateSong(new_song apitypes.ApiSong) (int, error) {
	user, status, err := users.GetUser(new_song.Token)
	if err != nil {
		return status, err
	}
	db := connection.GetDatabase()
	artist_username, status, err := new_song.Song.GetArtistUsername()
	if err != nil {
		return status, err
	}
	if artist_username != user.Username {
		return http.StatusUnauthorized, errors.New("User is not the owner of the song")
	}
	result := db.Collection("items").FindOne(context.TODO(), bson.D{{"_id", new_song.UUID}})
	var song apitypes.Song
	result.Decode(&song)
	song.Genre = new_song.Genre
	song.Price = new_song.Price
	song.Title = new_song.Title
	return http.StatusOK, nil
}
