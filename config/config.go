package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	DatabaseTimeout = 5 * time.Second
)

type Schema struct {
	Environment string `env:"environment"`
	Port        int    `env:"port"`
	AuthSecret  string `env:"auth_secret"`
	DatabaseURI string `env:"database_uri"`
}

var (
	ProductionEnv = "production"
	cfg           Schema
)

func init() {
	err := godotenv.Load("./config/config.yaml")
	if err != nil {
		log.Fatalf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}
}

func GetConfig() *Schema {
	return &cfg
}
