package users

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type Cart struct {
	Username string `bson:"_id"`
	Items    []apitypes.Buyable
}

func GetLoggedInUsersCart(token string) (Cart, int, error) {
	db := connection.GetDatabase()
	user, status, err := GetUser(token)
	if err != nil {
		return Cart{}, status, err
	}
	username := user.Username
	result := db.Collection("carts").FindOne(context.TODO(), bson.D{{"Username", username}})
	var cart Cart
	result.Decode(&cart)
	return cart, http.StatusOK, nil
}

func GetUsersCart(username string) (Cart, int, error) {
	db := connection.GetDatabase()
	result := db.Collection("carts").FindOne(context.TODO(), bson.D{{"Username", username}})
	var cart Cart
	result.Decode(&cart)
	return cart, http.StatusOK, nil
}

func AddToCart(token string, item apitypes.Buyable) (int, error) {
	db := connection.GetDatabase()
	cart, status, err := GetLoggedInUsersCart(token)
	if err != nil {
		return status, err
	}
	cart.Items = append(cart.Items, item)
	_, err = db.Collection("carts").UpdateByID(context.TODO(), cart.Username, cart)
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
		if v.UUID == item_uuid {
			to_remove_idx = i
			break
		}
	}
	cart.Items = append(cart.Items[:to_remove_idx], cart.Items[to_remove_idx+1:]...)
	_, err = db.Collection("carts").UpdateByID(context.TODO(), cart.Username, cart)
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
		total_cost += v.Price
	}
	if total_cost > user.Balance {
		return http.StatusBadRequest, errors.New("Not enough money.")
	}
	user.Balance -= total_cost
	cart.Items = cart.Items[:0]
	db.Collection("carts").UpdateByID(context.TODO(), cart.Username, cart)
	return http.StatusOK, nil
}
