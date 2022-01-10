package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Port string `yaml:"port"`
}

func (c *Config) LoadData() error {
	b, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, c)
	return err
}