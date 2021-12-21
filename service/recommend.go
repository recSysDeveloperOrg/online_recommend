package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	userObjectID         primitive.ObjectID
	totalRatingCnt       int64
	uninterestedMovieIds map[primitive.ObjectID]struct{}
	uninterestedTagIds   map[primitive.ObjectID]struct{}
	viewLogs             []primitive.ObjectID
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
}

func (*recommendService) checkParams(ctx *RecommendContext) {
	if len(strings.TrimSpace(ctx.req.UserId)) == 0 {
		ctx.errCode = BuildErrCode("没有用户ID信息", ParamErr)
		return
	}
	if ctx.req.Page < 0 {
		ctx.errCode = BuildErrCode(fmt.Sprintf("Page:%d", ctx.req.Page), ParamErr)
		return
	}
	if ctx.req.Offset < 0 {
		ctx.errCode = BuildErrCode(fmt.Sprintf("Offset:%d", ctx.req.Page), ParamErr)
	}
}

func (*recommendService) collectData(ctx *RecommendContext) {

}

func (*recommendService) buildResponse(ctx *RecommendContext) {
	errCode := Success
	if ctx.errCode != nil {
		errCode = ctx.errCode
	}
	ctx.resp.BaseResp.Code = errCode.Code
	ctx.resp.BaseResp.Msg = errCode.Msg
}
