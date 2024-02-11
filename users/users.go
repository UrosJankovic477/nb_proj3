package users

import (
	"context"
	"errors"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(username string, password string, displayname string) error {
	if len(displayname) == 0 {
		return errors.New("Display name cannot be empty")
	}

	if len(username) == 0 {
		return errors.New("Username cannot be empty")
	}

	username = strings.ToLower(username)

	if len(password) < 6 {
		return errors.New("Password must contain at least 6 characers")
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	db := connection.GetDatabase()
	_, err = db.Collection("users").InsertOne(context.TODO(), apitypes.User{
		Username:     username,
		PasswordHash: hash,
		DisplayName:  displayname,
		Balance:      100.0,
	})
	if err != nil {
		return err
	}
	cart := apitypes.Cart{Username: username, Items: make([]apitypes.IBuyable, 0)}
	db.Collection("carts").InsertOne(context.TODO(), cart)

	return nil
}

func Login(username string, password string) (string, int, error) {
	username = strings.ToLower(username)
	db := connection.GetDatabase()
	result := db.Collection("users").FindOne(context.TODO(), bson.D{{"_id", username}})
	var user apitypes.LoginUserTemp
	err := result.Decode(&user)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	correctPassword, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if !correctPassword {
		return "", http.StatusUnauthorized, errors.New("Wrong password")
	}

	sessionToken := GenerateToken(username)
	db.Collection("tokens").InsertOne(context.TODO(), sessionToken)
	return sessionToken.Token, http.StatusOK, nil
}

func GetLoggednInUsersData(token string) (apitypes.UserData, int, error) {
	user, status, err := GetUser(token)
	if err != nil {
		return apitypes.UserData{}, status, err
	}
	cart, _, _ := GetLoggedInUsersCart(token)

	userData := apitypes.UserData{
		DisplayName: user.DisplayName,
		Username:    user.Username,
		Balance:     user.Balance,
		CartCount:   uint16(len(cart.Items)),
	}
	return userData, http.StatusOK, nil
}

func GetUserFromUsername(username string) (apitypes.User, int, error) {
	db := connection.GetDatabase()
	result := db.Collection("users").FindOne(context.TODO(), bson.D{{"_id", username}})
	var user apitypes.User
	err := result.Decode(&user)
	if err != nil {
		return user, http.StatusNotFound, nil
	}
	return user, http.StatusOK, nil
}
