package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type MovieSimMatDao struct {
}

type MovieWeight struct {
	To     string  `bson:"from"`
	Weight float64 `bson:"weight"`
}
type MovieWeights struct {
	From string         `bson:"from"`
	To   []*MovieWeight `bson:"to"`
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
	c, err := GetClient().Collection(CollectionMovieSimMat).
		Find(ctx, condMap)
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
