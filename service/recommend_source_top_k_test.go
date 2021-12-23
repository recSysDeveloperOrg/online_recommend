package service

import (
	"context"
	"recommend/idl/gen/recommend"
	"testing"
)

func TestRecommendSourceTopK_RequestRecommend(t *testing.T) {
	svc := &RecommendSourceTopK{}
	rCtx := NewRecommendContext(context.Background(), &recommend.RecommendReq{
		UserId: testUserID,
		Page:   0,
		Offset: 10,
	})
	svc.RequestRecommend(rCtx)
	// should hit cache
	svc.RequestRecommend(rCtx)
	for _, recommendPairs := range rCtx.RecommendMovies {
		for _, recommendPair := range recommendPairs {
			t.Logf("%+v", recommendPair)
		}
	}
}
