package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AllowConfig struct {
	AllowedRegistries   []string `yaml:"AllowedRegistries"`
	AllowedRepositories []string `yaml:"AllowedRepositories"`
	AllowLatest         bool     `yaml:"AllowLatest"`
	AllowByDefault      bool     `yaml:"AllowByDefault"`
}

func (c *AllowConfig) ReadConfig(configPath string) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
