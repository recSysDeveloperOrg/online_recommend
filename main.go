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

	service.InitService()
}
