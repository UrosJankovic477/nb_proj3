package connection

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connection *mongo.Client = nil
var db *mongo.Database = nil

const uri = "mongodb://127.0.0.1:27017/"

func Init() error {
	if connection != nil && db != nil {
		return nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	connection, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	db = connection.Database("nb_proj3")
	return nil
}

func GetConnection() *mongo.Client {
	return connection
}

func GetDatabase() *mongo.Database {
	return db
}
