package service

import (
	"context"
	"recommend/constant"
	"recommend/idl/gen/recommend"
	"testing"
)

func TestRecommendService_DoService(t *testing.T) {
	rCtx := NewRecommendContext(context.Background(), &recommend.RecommendReq{
		UserId: testUserID,
		Page:   0,
		Offset: 10,
	})
	svc := NewRecommendService()
	svc.RecommendMovies(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetSuccess.Code {
		t.Fatal()
	}

	rCtx.Req.Page = -1
	svc.RecommendMovies(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetParamsErr.Code {
		t.Fatal()
	}

	rCtx.Req.Page = 0
	rCtx.Req.Offset = -1
	svc.RecommendMovies(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetParamsErr.Code {
		t.Fatal()
	}

	rCtx.Req.Page = 0
	rCtx.Req.Offset = 501
	svc.RecommendMovies(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetParamsErr.Code {
		t.Fatal()
	}
}
