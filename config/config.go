package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type LinkConfig struct {
	Expiration      uint   `yaml:"expiration"`
	TokenLength     int    `yaml:"token_length"`
	Alphabet        string `yaml:"alphabet"`
	RecreateRetries int    `yaml:"recreate_retries"`
}

type ServiceConfig struct {
	Host           string `yaml:"host"`
	Port           uint   `yaml:"port"`
	GrpcPort       uint   `yaml:"grpc_port"`
	RecalcInterval uint   `yaml:"recalculation_interval"`
}

type Config struct {
	LinkConfig    LinkConfig    `yaml:"link_config"`
	ServiceConfig ServiceConfig `yaml:"service_config"`
}

type DbConfig struct {
	Host             string `yaml:"host"`
	Port             uint   `yaml:"port"`
	Database         string `yaml:"database"`
	User             string `yaml:"user"`
	Password         string `yaml:"-"`
	ReconnectRetries int    `yaml:"reconnect_retries"`
}

type Summary struct {
	Cfg   Config   `yaml:"config"`
	DbCfg DbConfig `yaml:"db_config"`
}

func ParseConfig() (*Summary, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	filename, err := filepath.Abs("config/config.yml")

	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Summary{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	config.DbCfg.User = pgUser
	config.DbCfg.Password = pgPass
	return config, nil
}
