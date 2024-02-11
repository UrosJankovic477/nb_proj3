package search

import (
	"context"
	"nb_proj3/apitypes"
	"nb_proj3/connection"
	"nb_proj3/users"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchAlbums(query string, genre []string, max_price float64, count int64, page int64) ([]apitypes.ApiSearchResult, int, error) {
	db := connection.GetDatabase()
	opts := options.Find().SetSkip(page * count).SetLimit(count)
	//bsone := bson.E{}
	//if len(genre) > 0 {
	//	bsone = bson.E{Key: "genre", Value: bson.D{{"$in", genre}}}
	//}

	cursor, err := db.Collection("albums").Find(context.TODO(), bson.D{
		{"price", bson.D{
			{"$lte", max_price},
		}},
		{
			"$or", bson.A{
				bson.D{
					{"songs.genre", bson.M{
						"$in": genre,
					}},
				},
			},
		},
		{"$or", bson.A{
			bson.D{
				{"title", bson.D{
					{"$regex", query},
					{"$options", "i"},
				}},
				{"songs.title", bson.D{
					{"$regex", query},
					{"$options", "i"},
				}},
			},
		}},
	}, opts)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	albums := make([]apitypes.Album, 0)
	cursor.All(context.TODO(), &albums)
	items := make([]apitypes.ApiSearchResult, 0)
	for _, v := range albums {
		user, _, _ := users.GetUserFromUsername(v.ArtistUsername)

		items = append(items, apitypes.ApiSearchResult{
			UUID:              v.UUID,
			Title:             v.Title,
			ArtistUsername:    v.ArtistUsername,
			ArtistDisplayname: user.DisplayName,
			Price:             v.Price,
		})
	}

	return items, http.StatusOK, nil
}
