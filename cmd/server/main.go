package main

import (
	"fmt"
	"log"

	"github.com/fvbock/endless"
	"github.com/joho/godotenv"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	"github.com/thiri-lwin/gopher-tech-blog/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	cfg := config.LoadConfig()
	s := server.New(cfg)
	log.Println("Starting the server...")
	if err := endless.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), s); err != nil {
		panic(err)
	}
}
