package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type UserRatingDao struct {
}

var userRatingDao *UserRatingDao
var userRatingDaoOnce sync.Once

func NewUserRatingDao() *UserRatingDao {
	userRatingDaoOnce.Do(func() {
		userRatingDao = &UserRatingDao{}
	})

	return userRatingDao
}

type UserRating struct {
	UserID  string  `bson:"user_id"`
	MovieID string  `bson:"movie_id"`
	Rating  float64 `bson:"rating"`
}

func (*UserRatingDao) FindRatingAbove(ctx context.Context, userID string, minRating float64) ([]*UserRating, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	c, err := GetClient().Collection(CollectionRating).
		Find(ctx, bson.D{{"user_id", userObjectID}})
	if err != nil {
		return nil, err
	}

	var result []*UserRating
	if err := c.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
