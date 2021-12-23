package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type MovieModelDao struct {
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
	var res []string
	c, err := GetClient().Collection(CollectionMovie).
		Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"average_rating", -1}}).
			SetLimit(topK).SetProjection(bson.D{{"_id", 1}}))
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
