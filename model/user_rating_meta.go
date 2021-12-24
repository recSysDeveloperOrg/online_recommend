package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

type UserRatingMetaDao struct {
}

type UserRatingMeta struct {
	UserID      string `bson:"user_id"`
	TotalRating int64  `bson:"total_rating"`
}

var userRatingMetaDao *UserRatingMetaDao
var userRatingMetaDaoOnce sync.Once

func NewUserRatingMetaDao() *UserRatingMetaDao {
	userRatingMetaDaoOnce.Do(func() {
		userRatingMetaDao = &UserRatingMetaDao{}
	})

	return userRatingMetaDao
}

func (*UserRatingMetaDao) FindRatingCntByUserID(ctx context.Context, userID string) (int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}
	var res *UserRatingMeta
	if err := GetClient().Collection(CollectionUserRatingMeta).
		FindOne(ctx, bson.D{{"user_id", userObjectID}},
			options.FindOne()).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No rating for %s", userID)
			return 0, nil
		}
	}

	return res.TotalRating, nil
}
