package service

import (
	"recommend/cache"
	. "recommend/constant"
	. "recommend/idl/gen/recommend"
	"recommend/model"
	"sync"
)

type RecommendSourceTag struct {
}

var sourceTagUserID2RecommendPairCache = make(map[string][]*RecommendPair)

var recommendSourceTag *RecommendSourceTag
var recommendSourceTagOnce sync.Once

const (
	KMaxTag = 10 // 只使用该用户排名前10的标签进行推荐
)

func NewRecommendSourceTag() *RecommendSourceTag {
	recommendSourceTagOnce.Do(func() {
		recommendSourceTag = &RecommendSourceTag{}
	})

	return recommendSourceTag
}

func (*RecommendSourceTag) RequestRecommend(ctx *RecommendContext) {
	offset, size := ctx.Req.Page*ctx.Req.Offset, ctx.Req.Offset
	if recPair, hit := tryCache(sourceTagUserID2RecommendPairCache, ctx.Req.UserId, offset, size); hit {
		ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TAG] = recPair
		return
	}

	kMaxTags, err := model.NewTagUserDao().FindKMaxUserTags(ctx.Ctx,
		ctx.Req.UserId, KMaxTag)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	if len(kMaxTags) == 0 {
		return
	}

	tagID2Movies := make(map[string][]*model.TagMovie)
	for _, kMaxTag := range kMaxTags {
		kMaxTag2Movies, err := model.NewTagMovieDao().FindKMaxByTagID(ctx.Ctx,
			kMaxTag.TagID, MaxRecommend)
		if err != nil {
			ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
			return
		}
		tagID2Movies[kMaxTag.TagID] = kMaxTag2Movies
	}

	var heap *cache.Heap
	addedMovies := make(map[string]struct{})
	uninterestedSet, err := model.NewUserRecommendationMetaDao().FindUninterestedSet(ctx.Ctx, ctx.Req.UserId, model.UninterestedTypeTag)
	if err != nil {
		ctx.ErrCode = BuildErrCode(err, RetReadRepoErr)
		return
	}
	if uninterestedSet == nil {
		uninterestedSet = make(map[string]struct{})
	}

	for _, kMaxTag := range kMaxTags {
		if _, ok := uninterestedSet[kMaxTag.TagID]; ok {
			continue
		}
		kMaxTagMovies := tagID2Movies[kMaxTag.TagID]
		if heap == nil {
			initNodes := make([]*cache.HeapNode, len(kMaxTagMovies))
			for i, kMaxTagMovie := range kMaxTagMovies {
				initNodes[i] = &cache.HeapNode{
					Key: getTagWeight(len(kMaxTag.MovieIDs), kMaxTagMovie.TaggedTimes),
					Value: &RecommendPair{
						MovieID:  kMaxTagMovie.MovieID,
						SourceID: kMaxTag.TagID,
					},
				}
			}

			heap = cache.NewHeap(initNodes)
			continue
		}
		for _, kMaxTagMovie := range kMaxTagMovies {
			if _, ok := addedMovies[kMaxTagMovie.MovieID]; ok {
				continue
			}
			weight := getTagWeight(len(kMaxTag.MovieIDs), kMaxTagMovie.TaggedTimes)
			if weight.Compare(heap.TopKey()) > 0 {
				oldPair := heap.ReplaceTop(weight, &RecommendPair{
					MovieID:  kMaxTagMovie.MovieID,
					SourceID: kMaxTag.TagID,
				})
				addedMovies[kMaxTagMovie.MovieID] = struct{}{}
				delete(addedMovies, interface2RecommendPair(oldPair).MovieID)
			}
		}
	}

	if heap == nil {
		return
	}
	recommendPairs := interface2RecommendPairs(heap.PopValues())
	sourceTagUserID2RecommendPairCache[ctx.Req.UserId] = recommendPairs
	if offset >= int64(len(recommendPairs)) {
		return
	}
	rangeEnd := offset + size
	if offset + size > int64(len(recommendPairs)) {
		rangeEnd = int64(len(recommendPairs))
	}
	ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TAG] = recommendPairs[offset : rangeEnd]
}

func getTagWeight(userTagTimes int, movieTagTimes int64) float64Comparator {
	return float64Comparator(float64(userTagTimes) * float64(movieTagTimes))
}
