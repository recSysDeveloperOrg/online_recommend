package service

import (
	"recommend/config"
	"recommend/model"
	"testing"
)

const (
	testUserID  = "100000000000000000000001"
	testMovieID = "100000000000000000002571"
)

func TestMain(m *testing.M) {
	if err := config.InitConfig("../config/prod_conf.json"); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}
	InitService()
	if code := m.Run(); code != 0 {
		panic(code)
	}
}
