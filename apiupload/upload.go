package apiupload

import (
	"bytes"
	"errors"
	"nb_proj3/connection"
	"nb_proj3/makepreview"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UploadAudio(flie_path string, data []byte) (uint16, error) {
	content_type := http.DetectContentType(data)
	if !strings.HasPrefix(content_type, "audio/mpeg") {
		return 0, errors.New("invalid content type")
	}

	db := connection.GetDatabase()
	opts := options.GridFSBucket().SetName("audio")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(data)
	uploadOpts := options.GridFSUpload().SetMetadata(bson.D{{"type", "full_song"}})
	_, err = bucket.UploadFromStream(flie_path, reader, uploadOpts)
	if err != nil {
		return 0, err
	}
	preview, length, err := makepreview.ProcessSong(data)
	if err != nil {
		return 0, err
	}
	reader = bytes.NewReader(preview)
	uploadOpts = options.GridFSUpload().SetMetadata(bson.D{{"type", "preview"}})
	_, err = bucket.UploadFromStream("pre_"+flie_path, reader, uploadOpts)
	if err != nil {
		return 0, err
	}
	return length, nil
}

func UploadImage(album_uuid string, data []byte) error {
	content_type := http.DetectContentType(data)
	if !strings.HasPrefix(content_type, "image/") {
		return errors.New("invalid content type")
	}

	db := connection.GetDatabase()
	opts := options.GridFSBucket().SetName("images")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)
	_, err = bucket.UploadFromStream(album_uuid, reader)
	if err != nil {
		return err
	}

	return nil
}
