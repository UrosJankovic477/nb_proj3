package apiupload

import (
	"nb_proj3/connection"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteImage(content_id string) error {
	db := connection.GetDatabase()
	opts := options.GridFSBucket().SetName("images")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}
	id := findByFilename(content_id, "images")
	err = bucket.Delete(id)
	return err
}

func DeleteAudio(content_id string) error {
	db := connection.GetDatabase()
	opts := options.GridFSBucket().SetName("audio")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}
	id := findByFilename(content_id, "audio")
	err = bucket.Delete(id)
	return err
}
