package service

import (
	"context"
	"recommend/constant"
	"recommend/idl/gen/recommend"
	"testing"
)

func TestRecommendMetaService_AddFilterRule(t *testing.T) {
	svc := NewRecommendMetaService()
	rCtx := NewRecommendMetaContext(context.Background(), &recommend.FilterRuleReq{
		SourceId: testMovieID,
		FType:    recommend.FilterType_FILTER_TYPE_TAG,
		UserId:   testUserID,
	})
	svc.AddFilterRule(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetSuccess.Code {
		t.Fatal()
	}

	rCtx.Req.SourceId = ""
	svc.AddFilterRule(rCtx)
	if rCtx.Resp.BaseResp.Code != constant.RetParamsErr.Code {
		t.Fatal()
	}
}
