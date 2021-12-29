package service

import (
	. "recommend/constant"
	. "recommend/idl/gen/recommend"
	"recommend/model"
	"sync"
)

type RecommendSourceLog struct {
}

var sourceLogUserID2RecommendPairsCache = make(map[string][]*RecommendPair)

var recommendSourceLog *RecommendSourceLog
var recommendSourceLogOnce sync.Once

func NewRecommendSourceLog() *RecommendSourceLog {
	recommendSourceLogOnce.Do(func() {
		recommendSourceLog = &RecommendSourceLog{}
	})

	return recommendSourceLog
}

func (*RecommendSourceLog) RequestRecommend(ctx *RecommendContext) {
	offset, size := ctx.Req.Page*ctx.Req.Offset, ctx.Req.Offset
	if recPairs, cached := tryCache(sourceLogUserID2RecommendPairsCache, ctx.Req.UserId, offset, size); cached {
		ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_LOG] = recPairs
		return
	}

	movieIDs := model.NewUserRecommendationMetaDao().GetViewLog(ctx.Req.UserId)
	if len(movieIDs) == 0 {
		return
	}

	movieWeights, err := model.NewMovieSimMatDao().FindByMovieIDs(ctx.Ctx, movieIDs)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

	userRatings, err := model.NewUserRatingDao().FindRatingAbove(ctx.Ctx,
		ctx.Req.UserId, 0.0)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	ratedMovies := userRatings2MovieIDSet(userRatings)
	uninterestedMovies, err := model.NewUserRecommendationMetaDao().FindUninterestedSet(ctx.Ctx,
		ctx.Req.UserId, model.UninterestedTypeMovie)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

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
	sourceLogUserID2RecommendPairsCache[ctx.Req.UserId] = recommendPairs
	if offset >= int64(len(recommendPairs)) {
		return
	}
	rangeEnd := offset + size
	if offset + size > int64(len(recommendPairs)) {
		rangeEnd = int64(len(recommendPairs))
	}
	ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_LOG] = recommendPairs[offset : rangeEnd]
}
