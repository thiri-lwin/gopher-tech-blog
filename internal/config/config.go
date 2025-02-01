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
