package model

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"recommend/cache"
	"sync"
)

type UserRecommendationMetaDao struct {
	m *cache.LimitMap
}
type UserRecommendationMeta struct {
	UserID             string
	UninterestedMovies map[string]struct{}
	UninterestedTags   map[string]struct{}
	ViewLogs           []string
}

var userRecommendationMetaDao *UserRecommendationMetaDao
var userRecommendationMetaDaoOnce sync.Once

func NewUserRecommendationMetaDao() *UserRecommendationMetaDao {
	userRecommendationMetaDaoOnce.Do(func() {
		userRecommendationMetaDao = &UserRecommendationMetaDao{
			m: cache.NewLimitMap(100),
		}
	})

	return userRecommendationMetaDao
}

var (
	ErrNoElementUpdate = errors.New("update failed, no element got updated")
)

type UninterestedType string

const (
	UninterestedTypeMovie = "movie"
	UninterestedTypeTag   = "tag"

	UninterestedFieldName = "uninterested_%s_ids"
)

func (*UserRecommendationMetaDao) FindUninterestedSet(ctx context.Context, userID, typeName string) (
	map[string]struct{}, error) {
	userObjectId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var movies []string
	c, err := GetClient().Collection(CollectionUserRecommendationMeta).
		Find(ctx, bson.D{{"user_id", userObjectId}}, options.Find().SetProjection(
			bson.D{{fmt.Sprintf("uninterested_%s_ids", typeName), 1}}))
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &movies); err != nil {
		return nil, err
	}

	movieIDSet := make(map[string]struct{})
	for _, movieID := range movies {
		movieIDSet[movieID] = struct{}{}
	}

	return movieIDSet, nil
}

func (*UserRecommendationMetaDao) AddUninterestedSet(ctx context.Context, userID, itemID string,
	uType UninterestedType) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return err
	}

	whereMap := bson.D{
		{"user_id", userObjectID},
	}
	updatesMap := bson.D{
		{"$push", bson.D{
			{fmt.Sprintf(UninterestedFieldName, uType), itemObjectID},
		}},
	}
	updateRes, err := GetClient().Collection(CollectionUserRecommendationMeta).
		UpdateOne(ctx, whereMap, updatesMap)
	if err != nil {
		return err
	}
	if updateRes.UpsertedCount == 0 {
		return ErrNoElementUpdate
	}

	return nil
}

func (d *UserRecommendationMetaDao) AddViewLog(userID string, movieID string) error {
	d.m.Add(userID, movieID)
	return nil
}

func (d *UserRecommendationMetaDao) GetViewLog(userID string) []string {
	return d.m.GetItemsAsString(userID)
}
