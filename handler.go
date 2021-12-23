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
	rCtx := service.NewRecommendContext(ctx, req)
	service.NewRecommendService().RecommendMovies(rCtx)

	log.Printf("%+v", rCtx.Resp)
	return rCtx.Resp, nil
}

func (*RecommendServer) AddFilterRule(ctx context.Context, req *recommend.FilterRuleReq) (
	*recommend.FilterRuleResp, error) {
	log.Printf("%+v", req)
	rCtx := service.NewRecommendMetaContext(ctx, req)
	service.NewRecommendMetaService().AddFilterRule(rCtx)

	log.Printf("%+v", rCtx.Resp)
	return rCtx.Resp, nil
}
