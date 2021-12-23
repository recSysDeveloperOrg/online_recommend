package service

import (
	"context"
	. "recommend/constant"
	"recommend/idl/gen/recommend"
	"recommend/model"
	"strings"
	"sync"
)

type RecommendMetaService struct {
}

type RecommendMetaContext struct {
	Ctx     context.Context
	Req     *recommend.FilterRuleReq
	Resp    *recommend.FilterRuleResp
	ErrCode *ErrorCode
}

var recommendMetaService *RecommendMetaService
var recommendMetaServiceOnce sync.Once

var filterType2UninterestedType = map[recommend.FilterType]model.UninterestedType{
	recommend.FilterType_FILTER_TYPE_MOVIE: model.UninterestedTypeMovie,
	recommend.FilterType_FILTER_TYPE_TAG:   model.UninterestedTypeTag,
}

func NewRecommendMetaService() *RecommendMetaService {
	recommendMetaServiceOnce.Do(func() {
		recommendMetaService = &RecommendMetaService{}
	})

	return recommendMetaService
}

func NewRecommendMetaContext(ctx context.Context, req *recommend.FilterRuleReq) *RecommendMetaContext {
	return &RecommendMetaContext{
		Ctx: ctx,
		Req: req,
		Resp: &recommend.FilterRuleResp{
			BaseResp: &recommend.BaseResp{},
		},
	}
}

func (s *RecommendMetaService) AddFilterRule(ctx *RecommendMetaContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.ErrCode != nil {
		return
	}
	s.doService(ctx)
}

func (*RecommendMetaService) checkParams(ctx *RecommendMetaContext) {
	req := ctx.Req
	if strings.TrimSpace(req.SourceId) == "" {
		ctx.ErrCode = BuildErrCode("sourceID is empty", RetParamsErr)
		return
	}
}

func (*RecommendMetaService) doService(ctx *RecommendMetaContext) {
	if err := model.NewUserRecommendationMetaDao().AddUninterestedSet(ctx.Ctx, ctx.Req.UserId,
		ctx.Req.SourceId, filterType2UninterestedType[ctx.Req.FType]); err != nil {
		ctx.ErrCode = BuildErrCode(err, RetWriteRepoErr)
	}
}

func (*RecommendMetaService) buildResponse(ctx *RecommendMetaContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}

	ctx.Resp.BaseResp.Code = errCode.Code
	ctx.Resp.BaseResp.Msg = errCode.Msg
}
