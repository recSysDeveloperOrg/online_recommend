package service

import (
	. "recommend/constant"
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
	if r.isUserInColdStartState(ctx) {
		return
	}

	offset, size := ctx.req.Page*ctx.req.Offset, ctx.req.Offset
	if recPairs, cached := tryCache(sourceItemCFUserID2RecommendPairCache, ctx.req.UserId, offset, size); cached {
		ctx.recommendMovies[RecommendSourceTypeItemCF] = recPairs
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
	// 找不感兴趣的电影，通过评分函数过滤掉
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
	ctx.recommendMovies[RecommendSourceTypeItemCF] = recommendPairs[offset : offset+size]
}

func (*RecommendSourceSimMat) isUserInColdStartState(ctx *RecommendContext) bool {
	ratingCnt, err := model.NewUserRatingMetaDao().FindRatingCntByUserID(ctx.ctx, ctx.req.UserId)
	if err != nil {
		ctx.errCode = BuildErrCode(err, RetReadRepoErr)
		return true
	}

	return ratingCnt < ColdStartStateRatingCnt
}
