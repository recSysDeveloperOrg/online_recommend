package main

import (
	"context"
	"log"
	"recommend/idl/gen/recommend"
	"recommend/service"
)

type RecommendServer struct {
	*recommend.UnimplementedRecommenderServer
}

func (*RecommendServer) Recommend(ctx context.Context, req *recommend.RecommendReq) (
	*recommend.RecommendResp, error) {
	log.Printf("%+v", req)
	rCtx := &service.RecommendContext{
		Ctx: ctx,
		Req: req,
	}
	service.NewRecommendService().DoService(rCtx)

	log.Printf("%+v", rCtx.Resp)
	return rCtx.Resp, nil
}

func (*RecommendServer) AddFilterRule(ctx context.Context, req *recommend.FilterRuleReq) (
	*recommend.FilterRuleResp, error) {
	return nil, nil
}
