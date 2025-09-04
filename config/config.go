package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	CheckTimeMinute int    `env:"CHECK_TIME_MINUTES" env-required:"true"`
	RSSLinks        string `env:"RSS_LINKS" env-required:"true"`
	HomeServerURL   string `env:"HOME_SERVER_URL" env-required:"true"`
	Login           string `env:"LOGIN" env-required:"true"`
	Password        string `env:"PASSWORD" env-required:"true"`
}

const configFilePath = ".env"

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err == nil {
		log.Println("config variables loaded")
		return &cfg
	}

	log.Printf("error read config variables: %s ", err)

	log.Println("Trying to load from a .env file")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatal("config file .env does not exist")
	}

	err = cleanenv.ReadConfig(configFilePath, &cfg)
	if err != nil {
		log.Fatalf("error reading .env file: %s", err)
	}
	log.Println("config .env file loaded")
	return &cfg
}
