package apiupload

import (
	"context"
	"nb_proj3/connection"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findByFilename(filename string, bucket_name string) interface{} {
	db := connection.GetDatabase()
	var id interface{}

	id = nil

	opts := options.GridFSBucket().SetName(bucket_name)
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		return id
	}
	cursor, _ := bucket.Find(bson.D{{"filename", filename}})
	var results []gridfs.File

	cursor.All(context.TODO(), &results)
	for _, v := range results {
		id = v.ID
	}
	return id
}
