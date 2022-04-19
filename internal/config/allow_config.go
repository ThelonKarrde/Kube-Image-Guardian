package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AllowConfig struct {
	AllowedRegistries   []string          `yaml:"allowed_registries"`
	AllowedRepositories []string          `yaml:"allowed_repositories"`
	AllowLatest         bool              `yaml:"allow_latest"`
	AllowByDefault      bool              `yaml:"allow_by_default"`
	LogOnly             bool              `yaml:"log_only"`
	LimitsDefined       bool              `yaml:"limits_defined"`
	RequestsDefined     bool              `yaml:"requests_defined"`
	DesiredVersions     map[string]string `yaml:"desired_versions"`
}

func (c *AllowConfig) ReadConfig(configPath string) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("configFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(configFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
