package main

import (
	"log"

	"github.com/dspeirs7/animals/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	server.StartServer()
}
