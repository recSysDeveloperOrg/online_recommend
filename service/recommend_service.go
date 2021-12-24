package service

import (
	"context"
	"fmt"
	. "recommend/constant"
	"recommend/idl/gen/recommend"
	"sync"
)

type RecommendContext struct {
	Ctx     context.Context
	Req     *recommend.RecommendReq
	Resp    *recommend.RecommendResp
	ErrCode *ErrorCode

	RecommendMovies map[recommend.RecommendSourceType][]*RecommendPair
}

type recommendService struct {
}

var service *recommendService
var serviceOnce sync.Once

func NewRecommendService() *recommendService {
	serviceOnce.Do(func() {
		service = &recommendService{}
	})

	return service
}

func NewRecommendContext(ctx context.Context, req *recommend.RecommendReq) *RecommendContext {
	return &RecommendContext{
		Ctx: ctx,
		Req: req,
		Resp: &recommend.RecommendResp{
			BaseResp: &recommend.BaseResp{},
		},
		RecommendMovies: make(map[recommend.RecommendSourceType][]*RecommendPair),
	}
}

func (s *recommendService) RecommendMovies(ctx *RecommendContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.ErrCode != nil {
		return
	}
	s.doRecommend(ctx)
}

func (*recommendService) checkParams(ctx *RecommendContext) {
	if ctx.Req.Page < 0 {
		ctx.ErrCode = BuildErrCode(fmt.Sprintf("Page:%d", ctx.Req.Page), RetParamsErr)
		return
	}
	if ctx.Req.Offset < 0 {
		ctx.ErrCode = BuildErrCode(fmt.Sprintf("Offset:%d", ctx.Req.Page), RetParamsErr)
	}
	if ctx.Req.Page*ctx.Req.Offset+ctx.Req.Offset > MaxRecommend {
		ctx.ErrCode = BuildErrCode(fmt.Sprintf("max recommend is %d", MaxRecommend), RetParamsErr)
	}
}

func (*recommendService) doRecommend(ctx *RecommendContext) {
	// 未登录用户直接推荐TopK
	if ctx.Req.UserId == "" {
		recommendSourceTopK.RequestRecommend(ctx)
		return
	}

	for _, recommendSource := range RecommendSources {
		recommendSource.RequestRecommend(ctx)
		if len(ctx.RecommendMovies) > 0 {
			break
		}
	}
}

func (*recommendService) buildResponse(ctx *RecommendContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}
	ctx.Resp.BaseResp.Code = errCode.Code
	ctx.Resp.BaseResp.Msg = errCode.Msg

	for recType, recommendPairs := range ctx.RecommendMovies {
		for _, recommendPair := range recommendPairs {
			ctx.Resp.Entry = append(ctx.Resp.Entry, &recommend.RecommendEntry{
				RsType:   recType,
				MovieId:  recommendPair.MovieID,
				SourceId: recommendPair.SourceID,
			})
		}
	}
}
