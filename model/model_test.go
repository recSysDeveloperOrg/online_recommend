package model

import (
	"recommend/config"
	"testing"
)

func TestMain(m *testing.M) {
	if err := config.InitConfig("../config/prod_conf.json"); err != nil {
		panic(err)
	}
	if err := InitModel(); err != nil {
		panic(err)
	}
	if code := m.Run(); code != 0 {
		panic(code)
	}
}
