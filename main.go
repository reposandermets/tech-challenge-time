package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/reposandermets/tech-challenge-time/api"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failure reading .env %v", err)
	}
	api.Boot()
}
