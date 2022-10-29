package app

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml: "server"`
	Mongo  MongoConfig  `yaml: "mongo"`
}

type ServerConfig struct {
	Name string `yaml: "name"`
	Port *int64 `yaml: "port"`
}

type MongoConfig struct {
	URI        string `yaml: "uri"`
	Database   string `yaml: "database"`
	Username   string `yaml: "username"`
	Password   string `yaml: "password"`
	Collection string `yaml: "collection"`
}

func (c *Config) getConf(path string) *Config {

	// Read file yaml in configs folder
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func GetConfig(path string) *Config {
	var conf Config
	conf.getConf(path)
	server := conf.Server
	mongo := conf.Mongo

	return &Config{
		Server: server,
		Mongo:  mongo,
	}
}
