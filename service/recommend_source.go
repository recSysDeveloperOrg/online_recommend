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

// RatingFunc 评分计算函数，sourceID表示推荐源ID（电影/tag），targetID是推荐的电影ID，weight是初始推荐权重
type RatingFunc func(sourceID, targetID string, weight float64) float64

var RecommendSources []RecommendSource

const (
	MaxRecommend = 500 // 固定最大的推荐数量
)

func AppendRecommendSource(source ...RecommendSource) {
	RecommendSources = append(RecommendSources, source...)
}

func (f float64Comparator) Compare(comparator interface{}) int {
	if float64Comp, ok := comparator.(float64Comparator); ok {
		sub := float64(f) - float64(float64Comp)
		if sub > 0 {
			return 1
		} else if sub == 0 {
			return 0
		} else {
			return -1
		}
	}
	panic("compare failed")
}

func movieWeights2RecommendPairs(movieWeights []*model.MovieWeights, ratingFunc RatingFunc, limit int) []*RecommendPair {
	recommendPairs, nAdd := make([]*RecommendPair, limit), 0
	addedMovies := make(map[string]struct{})
	pointers := make([]int, len(movieWeights))
	for nAdd < limit {
		weight, sourceID, updated := maxMovieWeight(movieWeights, pointers, ratingFunc)
		if !updated {
			break
		}
		if weight == nil {
			continue
		}
		if _, ok := addedMovies[weight.To]; ok {
			continue
		}
		addedMovies[weight.To] = struct{}{}

		recommendPairs[nAdd] = &RecommendPair{
			MovieID:  weight.To,
			SourceID: sourceID,
		}
		nAdd++
	}

	return recommendPairs
}

func maxMovieWeight(movieWeights []*model.MovieWeights, pointers []int, ratingFunc RatingFunc) (*model.MovieWeight,
	string, bool) {
	maxPointer, maxWeight, maxSourceID := 0, float64(0), ""
	var maxMovieWeight *model.MovieWeight
	updated := false
	for i, pointer := range pointers {
		if pointer >= len(movieWeights[i].To) {
			continue
		}

		weight := movieWeights[i].To[pointer]
		if ratingFunc(movieWeights[i].From, weight.To, weight.Weight) > maxWeight {
			maxMovieWeight = movieWeights[i].To[pointer]
			maxPointer = i
			maxSourceID = movieWeights[i].From
		}
		updated = true
	}

	pointers[maxPointer]++
	if maxMovieWeight == nil {
		for i := 0; i < len(pointers); i++ {
			pointers[i]++
		}
	}
	return maxMovieWeight, maxSourceID, updated
}

// TODO sync this map 先别搞缓存，后面再说吧
func tryCache(cache map[string][]*RecommendPair, userID string, offset, size int64) ([]*RecommendPair, bool) {
	//if cachedPairs, ok := cache[userID]; ok {
	//	if offset+size <= int64(len(cachedPairs)) {
	//		return cachedPairs[offset : offset+size], true
	//	}
	//
	//	return nil, true
	//}

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
