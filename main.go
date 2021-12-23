package main

import (
	"recommend/config"
	"recommend/model"
	"recommend/service"
)

func main() {
	if err := config.InitConfig(config.DefaultCfg); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}

	rsItemCF := &service.RecommendSourceSimMat{}
	rsLog := &service.RecommendSourceLog{}
	rsTag := &service.RecommendSourceTag{}
	rsTopK := &service.RecommendSourceTopK{}
	service.AppendRecommendSource(rsItemCF, rsLog, rsTag, rsTopK)
}
