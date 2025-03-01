package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	defaultImgURL = "https://thiri-lwin.github.io/gopher-tech-blog-bucket"
)

var ImageURL string

type Config struct {
	Env         string `mapstructure:"ENV"`
	DatabaseURI string `mapstructure:"DB_URI"`
	Port        int    `mapstructure:"PORT"`
	PostLimit   int    `mapstructure:"POST_LIMIT"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`
	RedisUser   string `mapstructure:"REDIS_USER"`
	RedisPass   string `mapstructure:"REDIS_PASS"`
	ImageURL    string `mapstructure:"IMAGE_URL"`

	SMTPServer string `mapstructure:"SMTP_SERVER"`
	SMTPPort   int    `mapstructure:"SMTP_PORT"`

	EmailFrom         string `mapstructure:"EMAIL_FROM"`
	EmailPass         string `mapstructure:"EMAIL_PASS"`
	AdminEmail        string `mapstructure:"ADMIN_EMAIL"`
	JWTKey            string `mapstructure:"JWT_KEY"`
	JWTExpirationTime int    `mapstructure:"JWT_EXPIRATION_TIME"`
	GoogleClientID    string `mapstructure:"GOOGLE_CLIENT_ID"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv() // Automatically read environment variables

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if cfg.ImageURL == "" {
		cfg.ImageURL = defaultImgURL
	}
	ImageURL = cfg.ImageURL

	return &cfg
}
