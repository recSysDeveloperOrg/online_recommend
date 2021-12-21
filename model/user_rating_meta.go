package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type UserRatingMetaDao struct {
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
	var res int64
	if err := GetClient().Collection(CollectionUserRatingMeta).
		FindOne(ctx, bson.D{{"user_id", userObjectID}}).Decode(&res); err != nil {
		return 0, err
	}

	return res, nil
}
