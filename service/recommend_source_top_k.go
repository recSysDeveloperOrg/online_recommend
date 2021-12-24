package service

import (
	"context"
	. "recommend/idl/gen/recommend"
	"recommend/model"
	"sync"
)

type RecommendSourceTopK struct {
}

var topKMovieCacheSync = sync.RWMutex{}
var topKMovieCache []*RecommendPair

var recommendSourceTopK *RecommendSourceTopK
var recommendSourceTopKOnce sync.Once

func NewRecommendSourceTopK() *RecommendSourceTopK {
	recommendSourceTopKOnce.Do(func() {
		recommendSourceTopK = &RecommendSourceTopK{}
	})

	return recommendSourceTopK
}

func (*RecommendSourceTopK) RequestRecommend(ctx *RecommendContext) {
	topKMovieCacheSync.RLock()
	defer topKMovieCacheSync.RUnlock()
	offset, size := ctx.Req.Page*ctx.Req.Offset, ctx.Req.Offset
	ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TOP_K] = topKMovieCache[offset : offset+size]
}

// TODO 定时任务刷新
func (*RecommendSourceTopK) RefreshMovieCache() {
	topKMovieCacheSync.Lock()
	defer topKMovieCacheSync.Unlock()
	topKMovies, err := model.NewMovieDao().FindTopKMovies(context.Background(), MaxRecommend)
	if err != nil {
		panic("refresh movies failed")
	}
	topKMovieCache = make([]*RecommendPair, MaxRecommend)
	for i, topKMovie := range topKMovies {
		topKMovieCache[i] = &RecommendPair{
			MovieID: topKMovie,
		}
	}
}
