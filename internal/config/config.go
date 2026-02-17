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
	ServerModeReadOnly  = "ro"
	ServerModeReadWrite = "rw"

	defaultConfigPathLocal = "config/app.yml"
)

var GlobalCfg Config

func init() {
	if os.Getenv("GLOBAL_LOAD") == "false" {
		return
	}

	GlobalCfg = *MustLoad()
}

type Config struct {
	Service  ServiceConfig  `yaml:"service"`
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	SMTP     SMTPConfig     `yaml:"smtp"`
}

type ServerConfig struct {
	Port string `yaml:"port" env:"SERVER_PORT"`
	Mode string `yaml:"mode" env:"SERVER_MODE"`
}

type ServiceConfig struct {
	SessionTTL time.Duration `yaml:"session_ttl"`
}

type SMTPConfig struct {
	FromAddress string `yaml:"from_address" env:"GMAIL_ADDRESS"`
	Password    string `yaml:"password" env:"GMAIL_APP_PASSWORD"`
	Host        string `yaml:"host" env:"SMTP_HOST"`
	Port        string `yaml:"port" env:"SMTP_PORT"`
}

type PostgresConfig struct {
	Host         string `yaml:"host" env:"POSTGRES_HOST"`
	Port         string `yaml:"port" env:"POSTGRES_PORT"`
	User         string `yaml:"user" env:"POSTGRES_USER"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Database     string `yaml:"database" env:"POSTGRES_DATABASE"`
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
