package users

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetLoggedInUsersCart(token string) (apitypes.Cart, int, error) {
	db := connection.GetDatabase()
	user, status, err := GetUser(token)
	if err != nil {
		return apitypes.Cart{}, status, err
	}
	username := user.Username
	result := db.Collection("carts").FindOne(context.TODO(), bson.D{{"_id", username}})
	var temp_cart apitypes.TempCart
	result.Decode(&temp_cart)
	var cart apitypes.Cart
	cart.Username = temp_cart.Username
	cart.Items = make([]apitypes.IBuyable, 0)
	for _, v := range temp_cart.Items {
		item, status, err := apitypes.GetItem(v.UUID)
		if err != nil {
			return apitypes.Cart{}, status, err
		}
		cart.Items = append(cart.Items, item)
	}
	return cart, http.StatusOK, nil
}

func GetUsersCart(username string) (apitypes.Cart, int, error) {
	db := connection.GetDatabase()
	result := db.Collection("carts").FindOne(context.TODO(), bson.D{{"_id", username}})
	var temp_cart apitypes.TempCart
	result.Decode(&temp_cart)
	var cart apitypes.Cart
	cart.Username = temp_cart.Username
	cart.Items = make([]apitypes.IBuyable, 0)
	for _, v := range temp_cart.Items {
		item, _, _ := apitypes.GetItem(v.UUID)
		cart.Items = append(cart.Items, item)
	}
	return cart, http.StatusOK, nil
}

func AddSongToCart(token string, song_uuid string) (int, error) {
	db := connection.GetDatabase()
	cart, status, err := GetLoggedInUsersCart(token)
	if err != nil {
		return status, err
	}
	for _, v := range cart.Items {
		if v.GetUUID() == song_uuid {
			return http.StatusNotModified, nil
		}
	}
	result := db.Collection("songs").FindOne(context.TODO(), bson.D{{"_id", song_uuid}})
	var item apitypes.Song
	result.Decode(&item)
	cart.Items = append(cart.Items, item)
	_, err = db.Collection("carts").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func AddAlbumToCart(token string, album_uuid string) (int, error) {
	db := connection.GetDatabase()
	cart, status, err := GetLoggedInUsersCart(token)
	if err != nil {
		return status, err
	}
	for _, v := range cart.Items {
		if v.GetUUID() == album_uuid {
			return http.StatusNotModified, nil
		}
	}
	result := db.Collection("albums").FindOne(context.TODO(), bson.D{{"_id", album_uuid}})
	var item apitypes.Album
	result.Decode(&item)
	cart.Items = append(cart.Items, item)
	_, err = db.Collection("carts").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func RemoveFromCart(token string, item_uuid string) (int, error) {
	db := connection.GetDatabase()
	cart, status, err := GetLoggedInUsersCart(token)
	if err != nil {
		return status, err
	}
	var to_remove_idx int
	for i, v := range cart.Items {
		if v.GetUUID() == item_uuid {
			to_remove_idx = i
			break
		}
	}
	cart.Items = append(cart.Items[:to_remove_idx], cart.Items[to_remove_idx+1:]...)
	_, err = db.Collection("carts").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func Checkout(token string) (int, error) {
	db := connection.GetDatabase()
	user, status, err := GetUser(token)
	if err != nil {
		return status, err
	}
	cart, status, err := GetUsersCart(user.Username)
	if err != nil {
		return status, err
	}
	total_cost := 0.0
	for _, v := range cart.Items {
		price := v.GetPrice()
		total_cost += price
	}

	if total_cost > user.Balance {
		return http.StatusBadRequest, errors.New("Not enough money.")
	}

	user.Balance -= total_cost
	for _, v := range cart.Items {
		price := v.GetPrice()
		artist_username, status, err := v.GetArtistUsername()
		if err != nil {
			return status, err
		}
		artist, status, err := GetUserFromUsername(artist_username)
		if err != nil {
			return status, err
		}
		artist.Balance += price
		db.Collection("users").ReplaceOne(context.TODO(), bson.D{{"_id", artist_username}}, artist)
		user.PurchasedItems = append(user.PurchasedItems, v.GetUUID())
	}
	cart.Items = make([]apitypes.IBuyable, 0)
	db.Collection("users").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, user)
	db.Collection("carts").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, cart)
	return http.StatusOK, nil
}

func EmptyCart(token string) (int, error) {
	db := connection.GetDatabase()
	user, status, err := GetUser(token)
	if err != nil {
		return status, err
	}
	cart, status, err := GetUsersCart(user.Username)
	if err != nil {
		return status, err
	}

	cart.Items = make([]apitypes.IBuyable, 0)
	db.Collection("carts").ReplaceOne(context.TODO(), bson.D{{"_id", cart.Username}}, cart)
	return http.StatusOK, nil
}
