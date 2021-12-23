package service

import (
	"recommend/cache"
	. "recommend/idl/gen/recommend"
)

type RecommendSourceTag struct {
}

const (
	KMaxTag = 100 // 只使用该用户排名前100的标签进行推荐
)

var sourceTagUserID2RecommendPairCache = make(map[string][]*RecommendPair)

func (*RecommendSourceTag) RequestRecommend(ctx *RecommendContext) {
	offset, size := ctx.req.Page*ctx.req.Offset, ctx.req.Offset
	if recPair, hit := tryCache(sourceTagUserID2RecommendPairCache, ctx.req.UserId, offset, size); hit {
		ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TAG] = recPair
		return
	}

	var heap *cache.Heap
	addedMovies := make(map[string]struct{})
	for _, kMaxTag := range ctx.kMaxTags {
		kMaxTagMovies := ctx.kMaxTagID2Movies[kMaxTag.TagID]
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
			if _, ok := addedMovies[kMaxTagMovie.TagID]; ok {
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

	recommendPairs := interface2RecommendPairs(heap.PopValues())
	sourceTagUserID2RecommendPairCache[ctx.req.UserId] = recommendPairs
	ctx.recommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TAG] = recommendPairs[offset : offset+size]
}

func getTagWeight(userTagTimes int, movieTagTimes int64) float64Comparator {
	return float64Comparator(float64(userTagTimes) * float64(movieTagTimes))
}
