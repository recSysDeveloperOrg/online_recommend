package service

import (
	"context"
	"recommend/idl/gen/recommend"
	"testing"
)

func TestRecommendSourceTag_RequestRecommend(t *testing.T) {
	svc := &RecommendSourceTag{}
	rCtx := NewRecommendContext(context.Background(), &recommend.RecommendReq{
		UserId: "100000000000000000158665",
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
