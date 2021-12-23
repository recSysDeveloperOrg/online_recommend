package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type MovieModelDao struct {
}

type MovieModel struct {
	MovieID string `bson:"_id"`
}

var movieModelDao *MovieModelDao
var movieModelDaoOnce sync.Once

func NewMovieDao() *MovieModelDao {
	movieModelDaoOnce.Do(func() {
		movieModelDao = &MovieModelDao{}
	})

	return movieModelDao
}

func (*MovieModelDao) FindTopKMovies(ctx context.Context, topK int64) ([]string, error) {
	var res []*MovieModel
	c, err := GetClient().Collection(CollectionMovie).
		Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"average_rating", -1}}).
			SetLimit(topK).SetProjection(bson.D{{"_id", 1}}))
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &res); err != nil {
		return nil, err
	}

	movieIDs := make([]string, len(res))
	for i, movie := range res {
		movieIDs[i] = movie.MovieID
	}

	return movieIDs, nil
}
