package core

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var TICKERS []string
var KAFKA_HOST string
var KAFKA_PORT string

func Load() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	t := os.Getenv("TICKERS")
	TICKERS = strings.Split(t, ",")
	LoadTickers(TICKERS)

	KAFKA_HOST = os.Getenv("KAFKA_HOST")
	KAFKA_PORT = os.Getenv("KAFKA_PORT")
}
