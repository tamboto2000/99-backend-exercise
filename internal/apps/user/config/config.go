package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const cfgPath = "./configs/user_svc.yaml"

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type HTTP struct {
	Port string `yaml:"port"`
}

type Config struct {
	Database Database `yaml:"database"`
	HTTP     HTTP     `yaml:"http"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	f, err := os.Open(cfgPath)
	if err != nil {
		return cfg, err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)

	return cfg, err
}
