package service

import (
	"context"
	"recommend/idl/gen/recommend"
	"recommend/model"
	"testing"
)

func TestRecommendSourceLog_RequestRecommend(t *testing.T) {
	svc := &RecommendSourceLog{}
	_ = model.NewUserRecommendationMetaDao().AddViewLog(testUserID, testMovieID)
	rCtx := NewRecommendContext(context.Background(), &recommend.RecommendReq{
		UserId: testUserID,
		Page:   0,
		Offset: 20,
	})
	svc.RequestRecommend(rCtx)
	// hit cache
	svc.RequestRecommend(rCtx)
	for _, recommendPairs := range rCtx.RecommendMovies {
		for _, recommendPair := range recommendPairs {
			t.Logf("%+v", recommendPair)
		}
	}
}
