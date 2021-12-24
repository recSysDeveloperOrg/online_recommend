package service

import (
	"context"
	. "recommend/constant"
	"recommend/idl/gen/recommend"
	"recommend/model"
	"strings"
	"sync"
)

type ViewLogService struct {
}

type ViewLogContext struct {
	Ctx     context.Context
	Req     *recommend.ViewLogReq
	Resp    *recommend.ViewLogResp
	ErrCode *ErrorCode
}

var viewLogService *ViewLogService
var viewLogServiceOnce sync.Once

func NewViewLogService() *ViewLogService {
	viewLogServiceOnce.Do(func() {
		viewLogService = &ViewLogService{}
	})

	return viewLogService
}

func NewViewLogContext(ctx context.Context, req *recommend.ViewLogReq) *ViewLogContext {
	return &ViewLogContext{
		Ctx: ctx,
		Req: req,
		Resp: &recommend.ViewLogResp{
			BaseResp: &recommend.BaseResp{},
		},
	}
}

func (s *ViewLogService) AddViewLog(ctx *ViewLogContext) {
	defer s.buildResponse(ctx)
	if s.checkParam(ctx); ctx.ErrCode != nil {
		return
	}
	if s.addViewLog(ctx); ctx.ErrCode != nil {
		return
	}
}

func (*ViewLogService) checkParam(ctx *ViewLogContext) {
	req := ctx.Req
	if strings.TrimSpace(req.UserId) == "" {
		ctx.ErrCode = BuildErrCode("userID is empty", RetParamsErr)
		return
	}
	if strings.TrimSpace(req.MovieId) == "" {
		ctx.ErrCode = BuildErrCode("movieID is empty", RetParamsErr)
	}
}

func (*ViewLogService) addViewLog(ctx *ViewLogContext) {
	if err := model.NewUserRecommendationMetaDao().AddViewLog(ctx.Req.UserId, ctx.Req.MovieId); err != nil {
		ctx.ErrCode = BuildErrCode(err, RetWriteRepoErr)
	}
}

func (*ViewLogService) buildResponse(ctx *ViewLogContext) {
	errCode := RetSuccess
	if ctx.ErrCode != nil {
		errCode = ctx.ErrCode
	}

	ctx.Resp.BaseResp.Code = errCode.Code
	ctx.Resp.BaseResp.Msg = errCode.Msg
}
