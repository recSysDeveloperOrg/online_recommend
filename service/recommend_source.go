package service

import (
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

// RatingFunc 评分计算函数，sourceID表示推荐源ID（电影/tag），targetID是推荐的电影ID，weight是初始推荐权重
type RatingFunc func(sourceID, targetID string, weight float64) float64

const (
	RecommendSourceTypeItemCF = iota
	RecommendSourceTypeTag
	RecommendSourceTypeLog
	RecommendSourceTypeTopK

	MaxRecommend = 500 // 固定最大的推荐数量
)

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
	recommendPairs, nAdd := make([]*RecommendPair, 0, limit), 0
	addedMovies := make(map[string]struct{})
	pointers := make([]int, len(movieWeights))
	for nAdd < limit {
		weight, sourceID := maxMovieWeight(movieWeights, pointers, ratingFunc)
		if _, ok := addedMovies[weight.To]; ok {
			continue
		}

		recommendPairs[nAdd] = &RecommendPair{
			MovieID:  weight.To,
			SourceID: sourceID,
		}
	}

	return recommendPairs
}

func maxMovieWeight(movieWeights []*model.MovieWeights, pointers []int, ratingFunc RatingFunc) (*model.MovieWeight, string) {
	maxPointer, maxWeight, maxSourceID := 0, float64(0), ""
	var maxMovieWeight *model.MovieWeight
	for i, pointer := range pointers {
		weight := movieWeights[i].To[pointer]
		if ratingFunc(movieWeights[i].From, weight.To, weight.Weight) > maxWeight {
			maxMovieWeight = movieWeights[i].To[pointer]
			maxPointer = i
			maxSourceID = movieWeights[i].From
		}
	}

	pointers[maxPointer]++
	return maxMovieWeight, maxSourceID
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

func userRatings2MovieIDSet(ratings []*model.UserRating) map[string]struct{} {
	movieIDSet := make(map[string]struct{})
	for _, rating := range ratings {
		movieIDSet[rating.MovieID] = struct{}{}
	}

	return movieIDSet
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
