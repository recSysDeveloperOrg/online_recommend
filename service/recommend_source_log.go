package service

import (
	. "recommend/constant"
	"recommend/model"
)

type RecommendSourceLog struct {
}

var sourceLogUserID2RecommendPairsCache = make(map[string][]*RecommendPair)

func (*RecommendSourceLog) RequestRecommend(ctx *RecommendContext) {
	offset, size := ctx.req.Page*ctx.req.Offset, ctx.req.Offset
	if recPairs, cached := tryCache(sourceLogUserID2RecommendPairsCache, ctx.req.UserId, offset, size); cached {
		ctx.recommendMovies[RecommendSourceTypeLog] = recPairs
		return
	}

	movieIDs := ctx.viewLogs
	movieWeights, err := model.NewMovieSimMatDao().FindByMovieIDs(ctx.ctx, movieIDs)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

	recommendPairs := movieWeights2RecommendPairs(movieWeights, DefaultRatingFunc, MaxRecommend)
	sourceLogUserID2RecommendPairsCache[ctx.req.UserId] = recommendPairs
	ctx.recommendMovies[RecommendSourceTypeLog] = recommendPairs[offset : offset+size]
}
