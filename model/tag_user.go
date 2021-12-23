package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type TagUserDao struct {
}

type TagUser struct {
	UserID    string   `bson:"user_id"`
	UpdatedAt int64    `bson:"updated_at"`
	TagID     string   `bson:"tag_id"`
	MovieIDs  []string `bson:"movie_ids"`
}

var tagUserDao *TagUserDao
var tagUserDaoOnce sync.Once

func NewTagUserDao() *TagUserDao {
	tagUserDaoOnce.Do(func() {
		tagUserDao = &TagUserDao{}
	})

	return tagUserDao
}

func (*TagUserDao) FindKMaxUserTags(ctx context.Context, userID string, kMax int) ([]*TagUser, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var tagUsers []*TagUser
	match := bson.D{
		{"user_id", userObjectID},
	}
	calculateSize := bson.D{
		{"$addFields", bson.D{{"use_times", bson.D{{"$size", "movie_ids"}}}}}}
	sortBySizeDesc := bson.D{
		{"$sort", "use_times"},
	}
	limit := bson.D{
		{"$limit", kMax},
	}
	c, err := GetClient().Collection(CollectionTagUser).
		Aggregate(ctx, mongo.Pipeline{match, calculateSize, sortBySizeDesc, limit})
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &tagUsers); err != nil {
		return nil, err
	}

	return tagUsers, nil
}
