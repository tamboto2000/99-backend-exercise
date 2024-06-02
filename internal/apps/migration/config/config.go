package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const cfgPath = "./configs/migration.yaml"

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Service struct {
	Database      string `yaml:"database"`
	MigrationFile string `yaml:"migration_file"`
}

type Services map[string]Service

type Config struct {
	Database Database `yaml:"database"`
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
