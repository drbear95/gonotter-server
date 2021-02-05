package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	JWTRefreshSecret string `yaml:"jwt_refresh_secret"`
	MongoSecret string `yaml:"mongo_secret"`
}

func (c *Config) GetConfig() *Config {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}