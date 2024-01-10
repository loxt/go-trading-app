package main

import (
	"github.com/joho/godotenv"
	"github.com/loxt/go-trading-app/producer/internal/trades"
	"log"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	trades.LoadHostAndPort(os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))

	t := os.Getenv("TICKERS")

	topics := strings.Split(t, ",")

	for i, topic := range topics {
		topics[i] = strings.Trim(strings.Trim(topic, "\\"), "\"")
	}

	trades.SubscribeAndListen(topics)
}
