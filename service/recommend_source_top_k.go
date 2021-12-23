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

func (*RecommendSourceTopK) RequestRecommend(ctx *RecommendContext) {
	topKMovieCacheSync.RLock()
	defer topKMovieCacheSync.RUnlock()
	ctx.RecommendMovies[RecommendSourceType_RECOMMEND_SOURCE_TYPE_TOP_K] = topKMovieCache[:]
}

func (*RecommendSourceTopK) refreshMovieCache() {
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
