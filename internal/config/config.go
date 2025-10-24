package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	defaultConfigPathLocal = "config/local.yml"
)

type Config struct {
	Service  ServiceConfig  `yaml:"service"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServiceConfig struct {
	SessionTTL time.Duration `yaml:"session_ttl"`
}

type PostgresConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user" env:"POSTGRES_USER"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Database     string `yaml:"database"`
	Sslmode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file reading error:", err.Error())
	}

	_, err = os.Stat(defaultConfigPathLocal)
	if errors.Is(err, os.ErrNotExist) {
		panic("config file does not exist")
	}

	var cfg Config

	err = cleanenv.ReadConfig(defaultConfigPathLocal, &cfg)
	if err != nil {
		panic("error while reading config: " + err.Error())
	}

	return &cfg
}
