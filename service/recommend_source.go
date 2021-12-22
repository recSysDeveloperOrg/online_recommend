package service

import (
	"recommend/cache"
	"recommend/model"
)

type RecommendSource interface {
	RequestRecommend(ctx *RecommendContext)
}

type RecommendPair struct {
	MovieID  string
	SourceID string
}

type float64Comparator float64

type RecommendSourceType uint8

// RatingFunc 评分计算函数，sourceMovieID表示推荐源电影ID，在ITEM-CF中会用于和评分结合，weight是初始推荐权重
type RatingFunc func(sourceMovieID string, weight float64) float64Comparator

const (
	RecommendSourceTypeItemCF = iota
	RecommendSourceTypeTag
	RecommendSourceTypeLog
	RecommendSourceTypeTopK

	MaxRecommend = 500 // 固定最大的推荐数量
)

var DefaultRatingFunc = func(sourceMovieID string, weight float64) float64Comparator {
	return float64Comparator(weight)
}

func (f float64Comparator) Compare(comparator interface{}) int {
	if another, ok := comparator.(float64); ok {
		if float64(f) < another {
			return -1
		} else if float64(f) == another {
			return 0
		} else {
			return 1
		}
	}
	panic("compare failed")
}

func movieWeights2RecommendPairs(movieWeights []*model.MovieWeights, ratingFunc RatingFunc, limit int) []*RecommendPair {
	var heap *cache.Heap
	heapNodes := make([]*cache.HeapNode, 0, limit)
	addedMovies := make(map[string]interface{})
	for _, weights := range movieWeights {
		for _, weight := range weights.To {
			if len(heapNodes) < MaxRecommend {
				if _, ok := addedMovies[weight.To]; ok {
					continue
				}
				heapNodes = append(heapNodes, &cache.HeapNode{
					Key: ratingFunc(weights.From, weight.Weight),
					Value: &RecommendPair{
						MovieID:  weight.To,
						SourceID: weights.From,
					},
				})

				continue
			}
			if heap == nil {
				heap = cache.NewHeap(heapNodes)
			}
			if float64Comparator(weight.Weight).Compare(heap.TopKey()) > 0 {
				oldEntry := interface2RecommendPair(
					heap.ReplaceTop(float64Comparator(weight.Weight), &RecommendPair{
						MovieID:  weight.To,
						SourceID: weights.From,
					}))
				delete(addedMovies, oldEntry.MovieID)
				addedMovies[weight.To] = struct{}{}
			}
		}
	}

	return interface2RecommendPairs(heap.PopValues())
}

func tryCache(cache map[string][]*RecommendPair, userID string, offset, size int64) ([]*RecommendPair, bool) {
	if cachedPairs, ok := cache[userID]; ok {
		if offset+size < int64(len(cachedPairs)) {
			return cachedPairs[offset : offset+size], true
		}

		return nil, true
	}

	return nil, false
}

func interface2RecommendPairs(v []interface{}) []*RecommendPair {
	res := make([]*RecommendPair, len(v))
	for i := 0; i < len(v); i++ {
		res[i] = interface2RecommendPair(v[i])
	}

	return res
}

func interface2RecommendPair(v interface{}) *RecommendPair {
	if rp, ok := v.(*RecommendPair); ok {
		return rp
	}
	panic("nope")
}
