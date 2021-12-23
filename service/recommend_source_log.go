package service

import (
	. "recommend/constant"
	. "recommend/idl/gen/recommend"
	"recommend/model"
)

type RecommendSourceLog struct {
}

var sourceLogUserID2RecommendPairsCache = make(map[string][]*RecommendPair)

func (*RecommendSourceLog) RequestRecommend(ctx *RecommendContext) {
	offset, size := ctx.req.Page*ctx.req.Offset, ctx.req.Offset
	if recPairs, cached := tryCache(sourceLogUserID2RecommendPairsCache, ctx.req.UserId, offset, size); cached {
		ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_LOG] = recPairs
		return
	}

	movieIDs := ctx.viewLogs
	movieWeights, err := model.NewMovieSimMatDao().FindByMovieIDs(ctx.ctx, movieIDs)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

	uninterestedMovies, ratedMovies := ctx.uninterestedMovieIds, ctx.ratedMovies
	recommendPairs := movieWeights2RecommendPairs(movieWeights,
		func(sourceID, targetID string, weight float64) float64 {
			if _, ok := uninterestedMovies[targetID]; ok {
				return 0
			}
			if _, ok := ratedMovies[targetID]; ok {
				return 0
			}

			return weight
		}, MaxRecommend)
	sourceLogUserID2RecommendPairsCache[ctx.req.UserId] = recommendPairs
	ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_LOG] = recommendPairs[offset : offset+size]
}
