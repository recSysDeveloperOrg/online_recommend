package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Mongo *MongoDB       `json:"mongodb"`
	Es    *ElasticSearch `json:"elastic_search"`
}

type MongoDB struct {
	Url      string `json:"url"`
	DBName   string `json:"DBName"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ElasticSearch struct {
	Url   string `json:"url"`
	Index string `json:"index"`
}

var cfg Config

const (
	CfgFileMain   = "config/prod_%s_conf.json"
	CfgFileNested = "../config/prod_%s_conf.json"
)

func GetConfig() *Config {
	return &cfg
}

func InitConfig(cfgFile string) error {
	var env string
	flag.StringVar(&env, "env", "remote", "specify env")
	flag.Parse()
	log.Printf("env:%s", env)
	content, err := ioutil.ReadFile(fmt.Sprintf(cfgFile, env))
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, &cfg); err != nil {
		return err
	}

	return nil
}
