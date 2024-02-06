package users

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const expirationTime = 240

type Token struct {
	Username string
	Token    string
	Expires  int64
}

func GenerateToken(username string) Token {
	sessionToken := make([]byte, 32)
	io.ReadFull(rand.Reader, sessionToken)

	return Token{username, base64.StdEncoding.EncodeToString(sessionToken), time.Now().Unix() + int64(time.Hour.Seconds()*expirationTime)}
}

func GetUser(token string) (apitypes.User, int, error) {
	db := connection.GetDatabase()
	var found_token Token
	result := db.Collection("tokens").FindOne(context.TODO(), bson.D{{"Token", token}})
	err := result.Decode(&found_token)
	if err != nil {
		return apitypes.User{}, http.StatusInternalServerError, err
	}
	if found_token.Expires < time.Now().Unix() {
		db.Collection("tokens").DeleteOne(context.TODO(), bson.D{{"Token", token}})
		return apitypes.User{}, http.StatusUnauthorized, errors.New("Session expired. Please loggin again.")
	}
	var user apitypes.User
	result = db.Collection("users").FindOne(context.TODO(), bson.D{{"Username", found_token.Username}})
	err = result.Decode(&user)
	if err != nil {
		return apitypes.User{}, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}
