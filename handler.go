package main

import (
	"context"
	"recommend/idl/gen/recommend"
	"recommend/service"
)

type RecommendServer struct {
	*recommend.UnimplementedRecommenderServer
}

func (*RecommendServer) Recommend(ctx context.Context, req *recommend.RecommendReq) (
	*recommend.RecommendResp, error) {
	rCtx := &service.RecommendContext{
		Ctx: ctx,
		Req: req,
	}
	service.NewRecommendService().DoService(rCtx)

	return rCtx.Resp, nil
}

func (*RecommendServer) AddFilterRule(ctx context.Context, req *recommend.FilterRuleReq) (
	*recommend.FilterRuleResp, error) {
	return nil, nil
}

func (*RecommendServer) InvalidateSimMatCache(ctx context.Context, req *recommend.InvalidateSimMatCacheReq) (
	*recommend.InvalidateSimMatCacheResp, error) {
	return nil, nil
}
