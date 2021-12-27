package service

func InitService() {
	rsItemCF := NewRecommendSourceSimMat()
	rsLog := NewRecommendSourceLog()
	rsTag := NewRecommendSourceTag()
	rsTopK := NewRecommendSourceTopK()
	rsTopK.RefreshMovieCache()
	AppendRecommendSource(rsItemCF, rsTag, rsLog, rsTopK)
}
