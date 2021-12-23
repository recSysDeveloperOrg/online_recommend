package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type TagMovieDao struct {
}

type TagMovie struct {
	MovieID     string `bson:"movie_id"`
	TagID       string `bson:"tag_id"`
	UpdatedAt   int64  `bson:"updated_at"`
	TaggedTimes int64  `bson:"tagged_times"`
}

var tagMovieDao *TagMovieDao
var tagMovieDaoOnce sync.Once

func NewTagMovieDao() *TagMovieDao {
	tagMovieDaoOnce.Do(func() {
		tagMovieDao = &TagMovieDao{}
	})

	return tagMovieDao
}

func (*TagMovieDao) FindKMaxByTagID(ctx context.Context, tagID string, kMax int64) ([]*TagMovie, error) {
	tagObjectID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		return nil, err
	}

	c, err := GetClient().Collection(CollectionTagMovie).
		Find(ctx, bson.D{{"tag_id", tagObjectID}},
			options.Find().SetSort(bson.D{{"tagged_times", -1}}).SetLimit(kMax))
	if err != nil {
		return nil, err
	}

	var tagMovies []*TagMovie
	if err := c.All(ctx, &tagMovies); err != nil {
		return nil, err
	}

	return tagMovies, nil
}
