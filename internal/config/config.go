package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string `mapstructure:"ENV"`
	DatabaseURI string `mapstructure:"DB_URI"`
	Port        int    `mapstructure:"PORT"`
	PostLimit   int    `mapstructure:"POST_LIMIT"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`
	RedisUser   string `mapstructure:"REDIS_USER"`
	RedisPass   string `mapstructure:"REDIS_PASS"`
}

func LoadConfig() *Config {
	cfg := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &cfg
}
