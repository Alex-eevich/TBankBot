package config

import "os"

type Config struct {
	Token   string
	BaseURL string
}

func Load() *Config {
	return &Config{
		Token:   os.Getenv("TINKOFF_TOKEN"),
		BaseURL: "https://sandbox-invest-public-api.tbank.ru/rest",
	}
}
