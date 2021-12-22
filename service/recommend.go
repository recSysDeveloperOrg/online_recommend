package service

import (
	"context"
	"fmt"
	. "recommend/constant"
	"recommend/idl/gen/recommend"
	"strings"
	"sync"
)

type RecommendContext struct {
	ctx                  context.Context
	req                  *recommend.RecommendReq
	resp                 *recommend.RecommendResp
	errCode              *ErrorCode
	totalRatingCnt       int64
	uninterestedMovieIds map[string]struct{}
	uninterestedTagIds   map[string]struct{}
	viewLogs             []string

	recommendMovies map[RecommendSourceType][]*RecommendPair
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

func (s *recommendService) DoService(ctx *RecommendContext) {
	defer s.buildResponse(ctx)
	if s.checkParams(ctx); ctx.errCode != nil {
		return
	}
	s.doRecommend(ctx)
}

func (*recommendService) checkParams(ctx *RecommendContext) {
	if len(strings.TrimSpace(ctx.req.UserId)) == 0 {
		ctx.errCode = BuildErrCode("没有用户ID信息", RetParamsErr)
		return
	}
	if ctx.req.Page < 0 {
		ctx.errCode = BuildErrCode(fmt.Sprintf("Page:%d", ctx.req.Page), RetParamsErr)
		return
	}
	if ctx.req.Offset < 0 {
		ctx.errCode = BuildErrCode(fmt.Sprintf("Offset:%d", ctx.req.Page), RetParamsErr)
	}
}

func (*recommendService) doRecommend(ctx *RecommendContext) {

}

func (*recommendService) buildResponse(ctx *RecommendContext) {
	errCode := RetSuccess
	if ctx.errCode != nil {
		errCode = ctx.errCode
	}
	ctx.resp.BaseResp.Code = errCode.Code
	ctx.resp.BaseResp.Msg = errCode.Msg
}
