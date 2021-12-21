package main

import (
	"recommend/config"
	"recommend/model"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}

}
