package main

import (
	"fmt"
	"log"

	"github.com/fvbock/endless"
	"github.com/thiri-lwin/gopher-tech-blog/internal/server"
)

func main() {
	s := server.New()
	log.Println("Starting the server...")
	if err := endless.ListenAndServe(fmt.Sprintf(":%d", 8080), s); err != nil {
		panic(err)
	}
}
