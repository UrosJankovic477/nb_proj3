package apitypes

type ApiAlbum struct {
	Token string
	Album
}

type Buyable struct {
	UUID  string `bson:"_id"`
	Price float64
}

type Song struct {
	Title       string
	Link        string
	PreviewLink string
	Genre       string
	Length      uint32
	Buyable
}

type Album struct {
	Title          string
	ArtistUsername string
	AlbumArtLink   string
	Songs          []Song
	Buyable
}

type ApiAlbumRemove struct {
	Token string
	UUID  string
}

type User struct {
	Username     string `bson:"_id"`
	DisplayName  string
	PasswordHash string
	Balance      float64
	CartID       string
}

type UserData struct {
	Username  string
	Balance   float64
	CartCount uint16
}

type ApiAddToCart struct {
	Token string
	Buyable
}

type ApiRemoveFromCart struct {
	Token    string
	ItemUUID string
}
