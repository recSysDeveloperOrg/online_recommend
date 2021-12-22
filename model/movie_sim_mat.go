package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type MovieSimMatDao struct {
}

type MovieWeight struct {
	To     string
	Weight float64
}
type MovieWeights struct {
	From string
	To   []*MovieWeight
}

var movieSimMatDao *MovieSimMatDao
var movieSimMatDaoOnce sync.Once

func NewMovieSimMatDao() *MovieSimMatDao {
	movieSimMatDaoOnce.Do(func() {
		movieSimMatDao = &MovieSimMatDao{}
	})

	return movieSimMatDao
}

func (*MovieSimMatDao) FindByMovieIDs(ctx context.Context, movieIDs []string) ([]*MovieWeights, error) {
	movieObjectIDs := make([]primitive.ObjectID, len(movieIDs))
	for i, movieID := range movieIDs {
		movieObjectID, err := primitive.ObjectIDFromHex(movieID)
		if err != nil {
			return nil, err
		}
		movieObjectIDs[i] = movieObjectID
	}

	var result []*MovieWeights
	condMap := bson.D{{"from", bson.D{{"$in", movieObjectIDs}}}}
	projMap := bson.D{{"to", 1}}
	c, err := GetClient().Collection(CollectionMovieSimMat).
		Find(ctx, condMap, options.Find().SetProjection(projMap))
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
