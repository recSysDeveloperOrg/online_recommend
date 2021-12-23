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

	// 过滤不感兴趣的电影
	uninterestedMovies, err := model.NewUserRecommendationMetaDao().FindUninterestedSet(ctx.ctx,
		ctx.req.UserId, model.UninterestedTypeMovie)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	// 过滤掉已经评分过的电影
	userRatings, err := model.NewUserRatingDao().FindRatingAbove(ctx.ctx, ctx.req.UserId, 0.0)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	ratedMovies := userRatings2MovieIDSet(userRatings)
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
	ctx.recommendMovies[RecommendSourceTypeLog] = recommendPairs[offset : offset+size]
}
