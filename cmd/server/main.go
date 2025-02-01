package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fvbock/endless"
	"github.com/joho/godotenv"
	"github.com/thiri-lwin/gopher-tech-blog/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	dbURI := os.Getenv("DB_URI")
	port := os.Getenv("PORT")
	postInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to convert port to integer: %v", err)
	}
	s := server.New(dbURI)
	log.Println("Starting the server...")
	if err := endless.ListenAndServe(fmt.Sprintf(":%d", postInt), s); err != nil {
		panic(err)
	}
}
