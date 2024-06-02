package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const cfgPath = "./configs/public_api_svc.yaml"

type HTTP struct {
	Port string `yaml:"port"`
}

type Endpoints map[string]string

func (e Endpoints) Get(key string) string {
	return e[key]
}

type Service struct {
	Host      string    `yaml:"host"`
	Endpoints Endpoints `yaml:"endpoints"`
}

type Services map[string]Service

func (s Services) Get(key string) Service {
	return s[key]
}

type Config struct {
	HTTP     HTTP     `yaml:"http"`
	Services Services `yaml:"services"`
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
