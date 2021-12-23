package service

import (
	. "recommend/constant"
	. "recommend/idl/gen/recommend"
	"recommend/model"
)

type RecommendSourceSimMat struct {
}

const (
	MinRatingIfLike         = 3.0 // 最低的代表用户喜欢这个电影的评分
	ColdStartStateRatingCnt = 10  // 新用户必须打分这么多才能使用ITEM-CF算法进行推荐
)

var sourceItemCFUserID2RecommendPairCache = make(map[string][]*RecommendPair)

func (r *RecommendSourceSimMat) RequestRecommend(ctx *RecommendContext) {
	if ctx.totalRatingCnt < ColdStartStateRatingCnt {
		return
	}

	offset, size := ctx.req.Page*ctx.req.Offset, ctx.req.Offset
	if recPairs, cached := tryCache(sourceItemCFUserID2RecommendPairCache, ctx.req.UserId, offset, size); cached {
		ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_RATING] = recPairs
		return
	}

	ratings, err := model.NewUserRatingDao().FindRatingAbove(ctx.ctx, ctx.req.UserId, MinRatingIfLike)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

	movieIDs := make([]string, len(ratings))
	for i, rating := range ratings {
		movieIDs[i] = rating.MovieID
	}
	movieWeights, err := model.NewMovieSimMatDao().FindByMovieIDs(ctx.ctx, movieIDs)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return
	}

	movieID2Rating := make(map[string]float64)
	for _, rating := range ratings {
		movieID2Rating[rating.MovieID] = rating.Rating
	}
	// 找不感兴趣的电影，通过评分函数过滤掉 ; 过滤掉已经评分过的电影
	uninterestedMovies, ratedMovies := ctx.uninterestedMovieIds, ctx.ratedMovies
	recommendPairs := movieWeights2RecommendPairs(movieWeights,
		func(sourceMovieID, targetMovieID string, weight float64) float64 {
			if _, ok := uninterestedMovies[targetMovieID]; ok {
				return 0
			}
			if _, ok := ratedMovies[targetMovieID]; ok {
				return 0
			}

			return movieID2Rating[sourceMovieID] * weight
		}, MaxRecommend)
	sourceItemCFUserID2RecommendPairCache[ctx.req.UserId] = recommendPairs
	ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_RATING] = recommendPairs[offset : offset+size]
}
