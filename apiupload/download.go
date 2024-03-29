package apiupload

import (
	"io"
	"nb_proj3/connection"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DownloadImage(content_id string) (io.Reader, error) {
	db := connection.GetDatabase()
	opts := options.GridFSBucket().SetName("images")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}
	id := findByFilename(content_id, "images")
	stream, err := bucket.OpenDownloadStream(id)
	return stream, err
}

func DownloadAudio(content_id string) (io.Reader, error) {
	db := connection.GetDatabase()

	opts := options.GridFSBucket().SetName("audio")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}
	id := findByFilename(content_id, "audio")
	stream, err := bucket.OpenDownloadStream(id)
	return stream, err
}
