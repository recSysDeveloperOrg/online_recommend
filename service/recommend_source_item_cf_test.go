package service

import (
	"context"
	"recommend/idl/gen/recommend"
	"testing"
)

func TestRecommendSourceSimMat_RequestRecommend(t *testing.T) {
	svc := &RecommendSourceSimMat{}
	rCtx := NewRecommendContext(context.Background(), &recommend.RecommendReq{
		UserId: testUserID,
		Page:   0,
		Offset: 500,
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
