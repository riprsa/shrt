package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	DBName string `yaml:"dbname"`
	Args string `yaml:"args"`
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