package model

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

type UserRecommendationMetaDao struct {
	m map[string][]string
}

type UserRecommendationMeta struct {
	UserID             string   `bson:"user_id"`
	UninterestedMovies []string `bson:"uninterested_movie_ids"`
	UninterestedTags   []string `bson:"uninterested_tag_ids"`
	ViewLogs           []string `bson:"view_logs"`
}

var userRecommendationMetaDao *UserRecommendationMetaDao
var userRecommendationMetaDaoOnce sync.Once

func NewUserRecommendationMetaDao() *UserRecommendationMetaDao {
	userRecommendationMetaDaoOnce.Do(func() {
		userRecommendationMetaDao = &UserRecommendationMetaDao{
			m: make(map[string][]string),
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

	var meta UserRecommendationMeta
	if err := GetClient().Collection(CollectionUserRecommendationMeta).
		FindOne(ctx, bson.D{{"user_id", userObjectId}}, options.FindOne().SetProjection(
			bson.D{{fmt.Sprintf("uninterested_%s_ids", typeName), 1}})).Decode(&meta); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("no uninterested %s for userID %s", typeName, userID)
			return nil, nil
		}
	}

	idList := meta.UninterestedMovies
	if typeName == UninterestedTypeTag {
		idList = meta.UninterestedTags
	}
	idSet := make(map[string]struct{})
	for _, id := range idList {
		idSet[id] = struct{}{}
	}

	return idSet, nil
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
		UpdateOne(ctx, whereMap, updatesMap, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	if updateRes.UpsertedCount == 0 && updateRes.ModifiedCount == 0 {
		return ErrNoElementUpdate
	}

	return nil
}

func (d *UserRecommendationMetaDao) AddViewLog(userID string, movieID string) error {
	logs := d.m[userID]
	for _, logID := range logs {
		if logID == movieID {
			return nil
		}
	}

	if _, ok := d.m[userID]; !ok {
		d.m[userID] = make([]string, 0)
	}
	d.m[userID] = append(d.m[userID], movieID)
	return nil
}

func (d *UserRecommendationMetaDao) GetViewLog(userID string) []string {
	return d.m[userID]
}
