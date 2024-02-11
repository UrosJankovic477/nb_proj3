package apitypes

import (
	"context"
	"nb_proj3/connection"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pagination struct {
	Count int64
	Page  int64
}

type ApiGetAlbum struct {
	Token     string
	AlbumUUID string
}

type ApiAddAlbum struct {
	Token string
	Title string
}

type ApiAlbum struct {
	Token string
	Album
}

type IBuyable interface {
	GetPrice() float64
	GetUUID() string
	SetUUID(uuid string)
	GetArtistUsername() (string, int, error)
	GetTitle() string
	GetAlbumUUID() string
}

type Song struct {
	UUID      string `bson:"_id"`
	AlbumUUID string
	Title     string
	Genre     string
	Length    uint16
	Price     float64
}

func (song Song) GetPrice() float64 {
	return song.Price
}

func (song Song) GetUUID() string {
	return song.UUID
}

func (song Song) SetUUID(uuid string) {
	song.UUID = uuid
}

func (song Song) GetArtistUsername() (string, int, error) {
	db := connection.GetDatabase()
	result := db.Collection("songs").FindOne(context.TODO(), bson.D{{"_id", song.AlbumUUID}})
	var album Album
	err := result.Decode(&album)
	if err == mongo.ErrNoDocuments {
		return "", http.StatusNotFound, err
	}
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return album.ArtistUsername, http.StatusOK, nil
}

func (song Song) GetTitle() string {
	return song.Title
}

func (song Song) GetAlbumUUID() string {
	return song.AlbumUUID
}

type SongData struct {
	UUID  string
	Title string
	Price float64
	Genre string
}

type TempAlbums struct {
	Albums []Album
}

type Album struct {
	UUID           string `bson:"_id"`
	Title          string
	ArtistUsername string
	Price          float64
	Songs          []SongData
}

func (album Album) GetPrice() float64 {
	return album.Price
}

func (album Album) GetUUID() string {
	return album.UUID
}

func (album Album) SetUUID(uuid string) {
	album.UUID = uuid
}

func (album Album) GetArtistUsername() (string, int, error) {
	return album.ArtistUsername, http.StatusOK, nil
}
func (album Album) GetTitle() string {
	return album.Title
}

func (album Album) GetAlbumUUID() string {
	return album.UUID
}

func GetItem(item_uuid string) (IBuyable, int, error) {
	db := connection.GetDatabase()
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"_id", item_uuid}})
	var album Album
	err := result.Decode(&album)
	if err == mongo.ErrNoDocuments {
		result = db.Collection("songs").FindOne(context.TODO(), bson.D{{"_id", item_uuid}})
		var song Song
		err = result.Decode(&song)
		if err != nil {
			return Song{}, http.StatusInternalServerError, err
		}
		return song, http.StatusOK, nil
	}
	if err != nil {
		return Song{}, http.StatusInternalServerError, err
	}
	return album, http.StatusOK, nil
}

type AlbumData struct {
	UUID           string
	Title          string
	ArtistUsername string
	Songs          []Song
	Price          float64
	Created        bool
	Owns           bool
}

type ApiAlbumRemove struct {
	Token string
	UUID  string
}

type LoginUserTemp struct {
	Username     string `bson:"_id"`
	DisplayName  string
	PasswordHash string
	Balance      float64
}

type User struct {
	Username       string `bson:"_id"`
	DisplayName    string
	PasswordHash   string
	PurchasedItems []string
	Balance        float64
}

type UserData struct {
	Username    string
	DisplayName string
	Balance     float64
	CartCount   uint16
}

type TempCart struct {
	Username string `bson:"_id"`
	Items    []ApiCartItem
}

type Cart struct {
	Username string `bson:"_id"`
	Items    []IBuyable
}

type ApiAddToCart struct {
	Token    string
	ItemUUID string
}

type ApiRemoveFromCart struct {
	Token    string
	ItemUUID string
}

type ApiAddSongToAlbum struct {
	Token     string
	AlbumUUID string
	SongData
}

type ApiRemoveSongFromAlbum struct {
	Token    string
	SongUUID string
}

type ApiSong struct {
	Token string
	Song
}

type ApiSearch struct {
	Query    string
	MinPrice float64
	MaxPrice float64
	Genres   []string
	Pagination
}

type ApiSearchOptions struct {
	MaxPrice float64
	Genres   []string
}

type ApiSearchResult struct {
	UUID              string
	Title             string
	ArtistDisplayname string
	ArtistUsername    string
	Price             float64
}

type ApiCartItem struct {
	UUID      string  `bson:"_id"`
	AlbumUUID string  `bson:"albumuuid"`
	Title     string  `bson:"title"`
	Price     float64 `bson:"price"`
}
